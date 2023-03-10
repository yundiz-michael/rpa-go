package chromedp

import (
	"context"
	"fmt"
	"merkaba/chromedp/cdproto/cdp"
	"merkaba/chromedp/cdproto/css"
	"merkaba/chromedp/cdproto/dom"
	"merkaba/chromedp/cdproto/input"
	"merkaba/chromedp/cdproto/page"
	"merkaba/chromedp/cdproto/runtime"
	"time"
)

const (
	// ErrAttributeValueEmpty is element attribute is empty
	ErrAttributeValueEmpty   Error = "element attribute is empty"
	ErrAttributeSrcNoChanged Error = "element attribute src has not been changed"
	ErrHtmlNoChanged         Error = "element html has not been changed"
)

func AttrBy(sel interface{}, name string, value any, opts ...QueryOption) QueryAction {
	if value == nil {
		panic("value cannot be nil")
	}
	return JavascriptAttribute(sel, name, value, opts...)
}

// VisibleNodes is an element query action that retrieves the document element nodes
// matching the selector.
func VisibleNodes(sel interface{}, parentNode *cdp.Node, nodes *[]*cdp.Node, opts ...QueryOption) QueryAction {
	if nodes == nil {
		panic("nodes cannot be nil")
	}

	return QueryAfter(sel, nil, func(ctx context.Context, execCtx runtime.ExecutionContextID, n ...*cdp.Node) error {
		results := make([]*cdp.Node, 0)
		var clip page.Viewport
		var rect []float64
		for i, node := range n {
			if i == 0 && parentNode != nil {
				if err := CallFunctionOnNode(ctx, parentNode, getClientRectJS, &clip); err == nil {
					rect = make([]float64, 4)
					rect[0] = clip.X
					rect[1] = clip.Y
					rect[2] = clip.X + clip.Width
					rect[3] = clip.Y + clip.Height
				}
			}

			_, err := dom.GetBoxModel().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				continue
			}
			// check visibility
			var res bool
			err = CallFunctionOnNode(ctx, node, visibleJS, &res)
			if err != nil {
				continue
			}
			if !res {
				continue
			}
			if err := CallFunctionOnNode(ctx, node, getClientRectJS, &clip); err != nil {
				continue
			}
			/*不在可见区域，只计算y方向*/
			y := clip.Y + clip.Height
			if rect != nil && (y <= rect[1] || y >= rect[3]) {
				continue
			}
			results = append(results, node)
		}
		*nodes = results
		return nil
	}, opts...)
}

// ShowNode 显示一个隐藏的元素
func ShowNode(n *cdp.Node) Action {
	return ActionFunc(func(ctx context.Context) error {
		var res string
		err := CallFunctionOnNode(ctx, n, showJS, &res)
		if err != nil {
			return err
		}
		return nil
	})
}

/* 等待Image中SRC的属性已经READY */
func _waitImageSrcReady(s *Selector, param interface{}) {
	WaitFunc(s.WaitReady(func(ctx context.Context, execCtx runtime.ExecutionContextID, n *cdp.Node) error {
		n.RLock()
		defer n.RUnlock()
		value := n.AttributeValue("src")
		if len(value) > 0 {
			return nil
		}
		return ErrAttributeValueEmpty
	}))(s, param)
}

func WaitImageSrcReady(sel interface{}, mapValue map[string]string, opts ...QueryOption) QueryAction {
	return QueryAfter(sel, nil, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		for _, node := range nodes {
			mapValue[sel.(string)] = node.AttributeValue("src")
		}
		return nil
	}, append(opts, _waitImageSrcReady)...)
}

func WaitContentChanged(sel interface{}, param interface{}, opts ...QueryOption) QueryAction {
	return Query(sel, param, append(opts, _waitContentChanged)...)
}

func _waitContentChanged(s *Selector, param interface{}) {
	WaitFunc(s.WaitReady(func(ctx context.Context, execCtx runtime.ExecutionContextID, n *cdp.Node) error {
		n.RLock()
		defer n.RUnlock()
		if s.ActionTimeout > 0 {
			t := time.Now().UnixMilli() - s.createdTime
			if t >= s.ActionTimeout {
				return nil
			}
		}
		html, err := dom.GetOuterHTML().WithNodeID(n.NodeID).Do(ctx)
		if err != nil {
			return err
		}
		if html == param {
			return ErrHtmlNoChanged
		}
		return nil
	}))(s, param)
}

/* 等待Image中SRC的属性已经Change */
func _imageNodeSrcChanged(s *Selector, param interface{}) {
	WaitFunc(s.WaitReady(func(ctx context.Context, execCtx runtime.ExecutionContextID, n *cdp.Node) error {
		n.RLock()
		defer n.RUnlock()
		mapValue := param.(map[string]string)
		oldValue := mapValue[s.query.(string)]
		value := n.AttributeValue("src")
		if len(value) > 0 && value != oldValue {
			mapValue[s.query.(string)] = value
			return nil
		}
		return ErrAttributeSrcNoChanged
	}))(s, param)
}

func WaitImageSrcChanged(sel interface{}, values map[string]string, opts ...QueryOption) QueryAction {
	return Query(sel, values, append(opts, _imageNodeSrcChanged)...)
}

func ReadAttributes(sel interface{}, attrMap map[string]string, opts ...QueryOption) QueryAction {
	return QueryAfter(sel, nil, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) == 0 {
			return nil
		}
		node := nodes[0]
		computed, err := css.GetComputedStyleForNode(node.NodeID).Do(ctx)
		if err != nil {
			return err
		}
		for _, prop := range computed {
			name := prop.Name
			if _, ok := attrMap[name]; ok {
				attrMap[name] = prop.Value
			}
		}
		return nil
	}, append(opts, NodeVisible)...)
}

func MouseDragNode(n *cdp.Node, offsetX float64) ActionFunc {
	return func(ctx context.Context) error {
		var p *input.DispatchMouseEventParams
		t := cdp.ExecutorFromContext(ctx).(*Target)
		if t == nil {
			return ErrInvalidTarget
		}

		if err := dom.ScrollIntoViewIfNeeded().WithNodeID(n.NodeID).Do(ctx); err != nil {
			return err
		}

		boxes, err := dom.GetContentQuads().WithNodeID(n.NodeID).Do(ctx)
		if err != nil {
			return err
		}

		if len(boxes) == 0 {
			return ErrInvalidDimensions
		}

		box := boxes[0]
		var x, y float64
		c := len(box)
		if c%2 != 0 || c < 1 {
			return ErrInvalidDimensions
		}
		for i := 0; i < c; i += 2 {
			x += box[i]
			y += box[i+1]
		}
		x /= float64(c / 2)
		y /= float64(c / 2)

		p = &input.DispatchMouseEventParams{
			Type:       input.MousePressed,
			X:          x,
			Y:          y,
			Button:     input.Left,
			ClickCount: 1,
		}

		if err := p.Do(ctx); err != nil {
			return err
		}
		p.Type = input.MouseMoved
		tracks := BuildTracks(offsetX)
		for _, track := range tracks {
			p.X += track
			if err := p.Do(ctx); err != nil {
				return err
			}
		}
		time.Sleep(740 * time.Millisecond)
		p.Type = input.MouseReleased
		return p.Do(ctx)
	}
}

func MouseDrag(sel interface{}, offsetX float64, opts ...QueryOption) QueryAction {
	return QueryAfter(sel, offsetX, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) == 0 {
			return fmt.Errorf("找不到相关 Node")
		}
		return MouseDragNode(nodes[0], offsetX).Do(ctx)
	}, append(opts, NodeVisible)...)
}

func MouseOverNode(n *cdp.Node) ActionFunc {
	return func(ctx context.Context) error {
		var p *input.DispatchMouseEventParams
		t := cdp.ExecutorFromContext(ctx).(*Target)
		if t == nil {
			return ErrInvalidTarget
		}

		if err := dom.ScrollIntoViewIfNeeded().WithNodeID(n.NodeID).Do(ctx); err != nil {
			return err
		}

		boxes, err := dom.GetContentQuads().WithNodeID(n.NodeID).Do(ctx)
		if err != nil {
			return err
		}

		if len(boxes) == 0 {
			return ErrInvalidDimensions
		}

		box := boxes[0]
		var x, y float64
		c := len(box)
		if c%2 != 0 || c < 1 {
			return ErrInvalidDimensions
		}
		for i := 0; i < c; i += 2 {
			x += box[i]
			y += box[i+1]
		}
		x /= float64(c / 2)
		y /= float64(c / 2)

		p = &input.DispatchMouseEventParams{
			Type:       input.MouseMoved,
			X:          x,
			Y:          y,
			Button:     input.Left,
			ClickCount: 1,
		}
		return p.Do(ctx)
	}
}

func MouseOver(sel interface{}, opts ...QueryOption) QueryAction {
	return QueryAfter(sel, nil, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) == 0 {
			return fmt.Errorf("找不到相关 Node")
		}
		return MouseOverNode(nodes[0]).Do(ctx)
	}, append(opts, NodeVisible)...)
}
