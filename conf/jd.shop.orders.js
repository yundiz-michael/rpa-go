/*auth:@admin
  name:@jd.shop.orders
  desc:@文件描述
 */
function textOf(e){
    return e==null? "":e.text;
}

function parseSkuItem(_row,sku){
    _tds=_row.selects("td");
    if (_tds.length<7) return false;
    try{
        sku["id"]=_tds[0].text;
        sku["name"]=_tds[1].text;
        sku["price"]=parseFloat(_tds[2].text.replace("￥",""));
        sku["disaccount"]=parseFloat(_tds[3].text.replace("￥",""));
        sku["qty"]=parseInt(_tds[6].text);
    }catch{
        return false;
    }
    return true;
}

function parseOrderDetail(_page,result){
    /*物流信息*/
    var e=_page.select("span.wb-green",2);
    /*如果e是空，那么就是已取消状态wb-red*/
    if (e==null || e.text=="等待付款" ) {
        return false;
    }
    result["logisticsId"]=textOf(_page.select("#waybillid-1",1));
    result["logisticsName"]=textOf(_page.select("div.logistics-tab-left a.wb-blue",1));
    /*收货人信息*/
    var customer={};
    result["customer"]=customer;
    customer["dongdong"]=textOf(_page.select("//table[@id='receiveData']//tr[@name='dongdongICON']",1));
    rname=textOf(_page.select("//table[@id='receiveData']//tr[@class='copyable']"),1);
    customer["name"]=rname.replace("收货人:","");
    _page.click("#viewOrderMobile");
    page.wait(500);
    customer["mobile"]=textOf(_page.select("#mobile"));
    /*付款信息*/
    _tables=_page.selects("td.ordinf-td table");
    if (_tables.length>=3){
        _table=_tables[2];
        _trs=_table.selects("tr");
        for(var i=0;i<_trs.length;i++){
            _tr=_trs[i];
            ptext=_tr.text;
            if (ptext.startsWith("付款时间:")){
                result["paidTime"]=ptext.substring(6);
            }else if (ptext.startsWith("商品总额:")){
                result["amount"]=parseFloat(ptext.substring(6).replace("￥",""));
            }else if (ptext.startsWith("运费金额:")){
                result["freightAmount"]=parseFloat(ptext.substring(6).replace("￥",""));
            }else if (ptext.startsWith("促销价格:")){
                result["promotionPrice"]=parseFloat(ptext.substring(6).replace("￥",""));
            }else if (ptext.startsWith("优惠券:")){
                result["coupon"]=parseFloat(ptext.substring(5).replace("￥",""));
            }else if (ptext.startsWith("京豆:")){
                result["bean"]=parseFloat(ptext.substring(4).replace("￥",""));
            }else if (ptext.startsWith("应支付金额:")){
                result["paidAmount"]=parseFloat(ptext.substring(7).replace("￥",""));
            }
        }
        /*sku列表*/
        var items=[];
        result["items"]=items;
        _skuRows=_page.selects("div.mtb20 table.wb-table-b tr");
        for(i=1;i<_skuRows.length;i++){
            var oSkuItem={};
            if (!parseSkuItem(_skuRows[i],oSkuItem)) continue;
            items.push(oSkuItem);
        }
    }
    return true;
}

function parseOrder(_order,result){
    _id=_order.select("a.orderid-mr10");
    result["id"]=_id.text;
    _times=_order.selects("span.ml10");
    result["createdTime"]=_times[0].text.replace("下单时间：","");
    /*订单详情*/
    var detailPage=_id.clickPage();
    var state= parseOrderDetail(detailPage,result);
    detailPage.close();
    return state;
}

function isToday(str) {
    dateArray = str.split("-");
    var month = parseInt(dateArray[1]);
    var date = parseInt(dateArray[2]);
    var today = new Date();
    var b1 = month == (today.getMonth()+1);
    var b2 = date == today.getDate();
    return  b1 && b2;
}

/*如果今天还有数据返回true */
function parseToday(shopName,page){
    var _orders=page.selects("div.shopweb-table-wrapper");
    /*从1开始的原因是会命中标题*/
    var result=true;
    for(var i=1;i<_orders.length;i++){
        var orderData={"shopName":shopName};
        if (!parseOrder(_orders[i],orderData)) continue
        var createdTime = orderData["createdTime"];
        if (!isToday(createdTime)){
            console.log("今天最后一条数据日期:" + createdTime);
            result=false;
            break;
        }
        sendData(orderData);
    }
    return result;
}

exports.parseToday = parseToday;