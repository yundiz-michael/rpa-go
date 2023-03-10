// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package tracing

import (
	json "encoding/json"
	io "merkaba/chromedp/cdproto/io"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing(in *jlexer.Lexer, out *TraceConfig) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "recordMode":
			(out.RecordMode).UnmarshalEasyJSON(in)
		case "traceBufferSizeInKb":
			out.TraceBufferSizeInKb = float64(in.Float64())
		case "enableSampling":
			out.EnableSampling = bool(in.Bool())
		case "enableSystrace":
			out.EnableSystrace = bool(in.Bool())
		case "enableArgumentFilter":
			out.EnableArgumentFilter = bool(in.Bool())
		case "includedCategories":
			if in.IsNull() {
				in.Skip()
				out.IncludedCategories = nil
			} else {
				in.Delim('[')
				if out.IncludedCategories == nil {
					if !in.IsDelim(']') {
						out.IncludedCategories = make([]string, 0, 4)
					} else {
						out.IncludedCategories = []string{}
					}
				} else {
					out.IncludedCategories = (out.IncludedCategories)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.IncludedCategories = append(out.IncludedCategories, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "excludedCategories":
			if in.IsNull() {
				in.Skip()
				out.ExcludedCategories = nil
			} else {
				in.Delim('[')
				if out.ExcludedCategories == nil {
					if !in.IsDelim(']') {
						out.ExcludedCategories = make([]string, 0, 4)
					} else {
						out.ExcludedCategories = []string{}
					}
				} else {
					out.ExcludedCategories = (out.ExcludedCategories)[:0]
				}
				for !in.IsDelim(']') {
					var v2 string
					v2 = string(in.String())
					out.ExcludedCategories = append(out.ExcludedCategories, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "syntheticDelays":
			if in.IsNull() {
				in.Skip()
				out.SyntheticDelays = nil
			} else {
				in.Delim('[')
				if out.SyntheticDelays == nil {
					if !in.IsDelim(']') {
						out.SyntheticDelays = make([]string, 0, 4)
					} else {
						out.SyntheticDelays = []string{}
					}
				} else {
					out.SyntheticDelays = (out.SyntheticDelays)[:0]
				}
				for !in.IsDelim(']') {
					var v3 string
					v3 = string(in.String())
					out.SyntheticDelays = append(out.SyntheticDelays, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "memoryDumpConfig":
			if in.IsNull() {
				in.Skip()
				out.MemoryDumpConfig = nil
			} else {
				if out.MemoryDumpConfig == nil {
					out.MemoryDumpConfig = new(MemoryDumpConfig)
				}
				(*out.MemoryDumpConfig).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing(out *jwriter.Writer, in TraceConfig) {
	out.RawByte('{')
	first := true
	_ = first
	if in.RecordMode != "" {
		const prefix string = ",\"recordMode\":"
		first = false
		out.RawString(prefix[1:])
		(in.RecordMode).MarshalEasyJSON(out)
	}
	if in.TraceBufferSizeInKb != 0 {
		const prefix string = ",\"traceBufferSizeInKb\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float64(float64(in.TraceBufferSizeInKb))
	}
	if in.EnableSampling {
		const prefix string = ",\"enableSampling\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.EnableSampling))
	}
	if in.EnableSystrace {
		const prefix string = ",\"enableSystrace\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.EnableSystrace))
	}
	if in.EnableArgumentFilter {
		const prefix string = ",\"enableArgumentFilter\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.EnableArgumentFilter))
	}
	if len(in.IncludedCategories) != 0 {
		const prefix string = ",\"includedCategories\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v4, v5 := range in.IncludedCategories {
				if v4 > 0 {
					out.RawByte(',')
				}
				out.String(string(v5))
			}
			out.RawByte(']')
		}
	}
	if len(in.ExcludedCategories) != 0 {
		const prefix string = ",\"excludedCategories\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v6, v7 := range in.ExcludedCategories {
				if v6 > 0 {
					out.RawByte(',')
				}
				out.String(string(v7))
			}
			out.RawByte(']')
		}
	}
	if len(in.SyntheticDelays) != 0 {
		const prefix string = ",\"syntheticDelays\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v8, v9 := range in.SyntheticDelays {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	if in.MemoryDumpConfig != nil {
		const prefix string = ",\"memoryDumpConfig\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.MemoryDumpConfig).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v TraceConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TraceConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TraceConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TraceConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing1(in *jlexer.Lexer, out *StartParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "bufferUsageReportingInterval":
			out.BufferUsageReportingInterval = float64(in.Float64())
		case "transferMode":
			(out.TransferMode).UnmarshalEasyJSON(in)
		case "streamFormat":
			(out.StreamFormat).UnmarshalEasyJSON(in)
		case "streamCompression":
			(out.StreamCompression).UnmarshalEasyJSON(in)
		case "traceConfig":
			if in.IsNull() {
				in.Skip()
				out.TraceConfig = nil
			} else {
				if out.TraceConfig == nil {
					out.TraceConfig = new(TraceConfig)
				}
				(*out.TraceConfig).UnmarshalEasyJSON(in)
			}
		case "perfettoConfig":
			out.PerfettoConfig = string(in.String())
		case "tracingBackend":
			(out.TracingBackend).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing1(out *jwriter.Writer, in StartParams) {
	out.RawByte('{')
	first := true
	_ = first
	if in.BufferUsageReportingInterval != 0 {
		const prefix string = ",\"bufferUsageReportingInterval\":"
		first = false
		out.RawString(prefix[1:])
		out.Float64(float64(in.BufferUsageReportingInterval))
	}
	if in.TransferMode != "" {
		const prefix string = ",\"transferMode\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.TransferMode).MarshalEasyJSON(out)
	}
	if in.StreamFormat != "" {
		const prefix string = ",\"streamFormat\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.StreamFormat).MarshalEasyJSON(out)
	}
	if in.StreamCompression != "" {
		const prefix string = ",\"streamCompression\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.StreamCompression).MarshalEasyJSON(out)
	}
	if in.TraceConfig != nil {
		const prefix string = ",\"traceConfig\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.TraceConfig).MarshalEasyJSON(out)
	}
	if in.PerfettoConfig != "" {
		const prefix string = ",\"perfettoConfig\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PerfettoConfig))
	}
	if in.TracingBackend != "" {
		const prefix string = ",\"tracingBackend\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.TracingBackend).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v StartParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v StartParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *StartParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *StartParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing1(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing2(in *jlexer.Lexer, out *RequestMemoryDumpReturns) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "dumpGuid":
			out.DumpGUID = string(in.String())
		case "success":
			out.Success = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing2(out *jwriter.Writer, in RequestMemoryDumpReturns) {
	out.RawByte('{')
	first := true
	_ = first
	if in.DumpGUID != "" {
		const prefix string = ",\"dumpGuid\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.DumpGUID))
	}
	if in.Success {
		const prefix string = ",\"success\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Success))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestMemoryDumpReturns) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestMemoryDumpReturns) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestMemoryDumpReturns) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestMemoryDumpReturns) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing2(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing3(in *jlexer.Lexer, out *RequestMemoryDumpParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "deterministic":
			out.Deterministic = bool(in.Bool())
		case "levelOfDetail":
			(out.LevelOfDetail).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing3(out *jwriter.Writer, in RequestMemoryDumpParams) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Deterministic {
		const prefix string = ",\"deterministic\":"
		first = false
		out.RawString(prefix[1:])
		out.Bool(bool(in.Deterministic))
	}
	if in.LevelOfDetail != "" {
		const prefix string = ",\"levelOfDetail\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.LevelOfDetail).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestMemoryDumpParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestMemoryDumpParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestMemoryDumpParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestMemoryDumpParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing3(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing4(in *jlexer.Lexer, out *RecordClockSyncMarkerParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "syncId":
			out.SyncID = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing4(out *jwriter.Writer, in RecordClockSyncMarkerParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"syncId\":"
		out.RawString(prefix[1:])
		out.String(string(in.SyncID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RecordClockSyncMarkerParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RecordClockSyncMarkerParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RecordClockSyncMarkerParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RecordClockSyncMarkerParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing4(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing5(in *jlexer.Lexer, out *MemoryDumpConfig) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing5(out *jwriter.Writer, in MemoryDumpConfig) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MemoryDumpConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MemoryDumpConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MemoryDumpConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MemoryDumpConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing5(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing6(in *jlexer.Lexer, out *GetCategoriesReturns) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "categories":
			if in.IsNull() {
				in.Skip()
				out.Categories = nil
			} else {
				in.Delim('[')
				if out.Categories == nil {
					if !in.IsDelim(']') {
						out.Categories = make([]string, 0, 4)
					} else {
						out.Categories = []string{}
					}
				} else {
					out.Categories = (out.Categories)[:0]
				}
				for !in.IsDelim(']') {
					var v10 string
					v10 = string(in.String())
					out.Categories = append(out.Categories, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing6(out *jwriter.Writer, in GetCategoriesReturns) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.Categories) != 0 {
		const prefix string = ",\"categories\":"
		first = false
		out.RawString(prefix[1:])
		{
			out.RawByte('[')
			for v11, v12 := range in.Categories {
				if v11 > 0 {
					out.RawByte(',')
				}
				out.String(string(v12))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetCategoriesReturns) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetCategoriesReturns) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetCategoriesReturns) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetCategoriesReturns) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing6(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing7(in *jlexer.Lexer, out *GetCategoriesParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing7(out *jwriter.Writer, in GetCategoriesParams) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetCategoriesParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetCategoriesParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetCategoriesParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetCategoriesParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing7(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing8(in *jlexer.Lexer, out *EventTracingComplete) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "dataLossOccurred":
			out.DataLossOccurred = bool(in.Bool())
		case "stream":
			out.Stream = io.StreamHandle(in.String())
		case "traceFormat":
			(out.TraceFormat).UnmarshalEasyJSON(in)
		case "streamCompression":
			(out.StreamCompression).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing8(out *jwriter.Writer, in EventTracingComplete) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"dataLossOccurred\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.DataLossOccurred))
	}
	if in.Stream != "" {
		const prefix string = ",\"stream\":"
		out.RawString(prefix)
		out.String(string(in.Stream))
	}
	if in.TraceFormat != "" {
		const prefix string = ",\"traceFormat\":"
		out.RawString(prefix)
		(in.TraceFormat).MarshalEasyJSON(out)
	}
	if in.StreamCompression != "" {
		const prefix string = ",\"streamCompression\":"
		out.RawString(prefix)
		(in.StreamCompression).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventTracingComplete) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventTracingComplete) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventTracingComplete) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventTracingComplete) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing8(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing9(in *jlexer.Lexer, out *EventDataCollected) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "value":
			if in.IsNull() {
				in.Skip()
				out.Value = nil
			} else {
				in.Delim('[')
				if out.Value == nil {
					if !in.IsDelim(']') {
						out.Value = make([]easyjson.RawMessage, 0, 2)
					} else {
						out.Value = []easyjson.RawMessage{}
					}
				} else {
					out.Value = (out.Value)[:0]
				}
				for !in.IsDelim(']') {
					var v13 easyjson.RawMessage
					(v13).UnmarshalEasyJSON(in)
					out.Value = append(out.Value, v13)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing9(out *jwriter.Writer, in EventDataCollected) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"value\":"
		out.RawString(prefix[1:])
		if in.Value == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v14, v15 := range in.Value {
				if v14 > 0 {
					out.RawByte(',')
				}
				(v15).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventDataCollected) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventDataCollected) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventDataCollected) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventDataCollected) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing9(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing10(in *jlexer.Lexer, out *EventBufferUsage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "percentFull":
			out.PercentFull = float64(in.Float64())
		case "eventCount":
			out.EventCount = float64(in.Float64())
		case "value":
			out.Value = float64(in.Float64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing10(out *jwriter.Writer, in EventBufferUsage) {
	out.RawByte('{')
	first := true
	_ = first
	if in.PercentFull != 0 {
		const prefix string = ",\"percentFull\":"
		first = false
		out.RawString(prefix[1:])
		out.Float64(float64(in.PercentFull))
	}
	if in.EventCount != 0 {
		const prefix string = ",\"eventCount\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float64(float64(in.EventCount))
	}
	if in.Value != 0 {
		const prefix string = ",\"value\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float64(float64(in.Value))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventBufferUsage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventBufferUsage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventBufferUsage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventBufferUsage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing10(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing11(in *jlexer.Lexer, out *EndParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing11(out *jwriter.Writer, in EndParams) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EndParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing11(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EndParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoTracing11(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EndParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing11(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EndParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoTracing11(l, v)
}