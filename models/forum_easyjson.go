// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
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

func easyjsonC8d74561DecodeGithubComNdRTechDbForumModels(in *jlexer.Lexer, out *TreadArr) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(TreadArr, 0, 1)
			} else {
				*out = TreadArr{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Thread
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC8d74561EncodeGithubComNdRTechDbForumModels(out *jwriter.Writer, in TreadArr) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v TreadArr) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8d74561EncodeGithubComNdRTechDbForumModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TreadArr) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8d74561EncodeGithubComNdRTechDbForumModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TreadArr) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8d74561DecodeGithubComNdRTechDbForumModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TreadArr) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8d74561DecodeGithubComNdRTechDbForumModels(l, v)
}
func easyjsonC8d74561DecodeGithubComNdRTechDbForumModels1(in *jlexer.Lexer, out *ThreadUpdate) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "message":
			out.Message = string(in.String())
		case "title":
			out.Title = string(in.String())
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
func easyjsonC8d74561EncodeGithubComNdRTechDbForumModels1(out *jwriter.Writer, in ThreadUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"message\":")
	out.String(string(in.Message))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"title\":")
	out.String(string(in.Title))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8d74561EncodeGithubComNdRTechDbForumModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8d74561EncodeGithubComNdRTechDbForumModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8d74561DecodeGithubComNdRTechDbForumModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8d74561DecodeGithubComNdRTechDbForumModels1(l, v)
}
func easyjsonC8d74561DecodeGithubComNdRTechDbForumModels2(in *jlexer.Lexer, out *Thread) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			if in.IsNull() {
				in.Skip()
				out.Id = nil
			} else {
				if out.Id == nil {
					out.Id = new(int)
				}
				*out.Id = int(in.Int())
			}
		case "slug":
			if in.IsNull() {
				in.Skip()
				out.Slug = nil
			} else {
				if out.Slug == nil {
					out.Slug = new(string)
				}
				*out.Slug = string(in.String())
			}
		case "title":
			out.Title = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "forum":
			out.Forum_title = string(in.String())
		case "author":
			out.User_nick = string(in.String())
		case "created":
			if in.IsNull() {
				in.Skip()
				out.Created = nil
			} else {
				if out.Created == nil {
					out.Created = new(string)
				}
				*out.Created = string(in.String())
			}
		case "votes":
			if in.IsNull() {
				in.Skip()
				out.Votes_count = nil
			} else {
				if out.Votes_count == nil {
					out.Votes_count = new(int)
				}
				*out.Votes_count = int(in.Int())
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
func easyjsonC8d74561EncodeGithubComNdRTechDbForumModels2(out *jwriter.Writer, in Thread) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	if in.Id == nil {
		out.RawString("null")
	} else {
		out.Int(int(*in.Id))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"slug\":")
	if in.Slug == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Slug))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"title\":")
	out.String(string(in.Title))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"message\":")
	out.String(string(in.Message))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"forum\":")
	out.String(string(in.Forum_title))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"author\":")
	out.String(string(in.User_nick))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"created\":")
	if in.Created == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Created))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"votes\":")
	if in.Votes_count == nil {
		out.RawString("null")
	} else {
		out.Int(int(*in.Votes_count))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Thread) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8d74561EncodeGithubComNdRTechDbForumModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Thread) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8d74561EncodeGithubComNdRTechDbForumModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Thread) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8d74561DecodeGithubComNdRTechDbForumModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Thread) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8d74561DecodeGithubComNdRTechDbForumModels2(l, v)
}
func easyjsonC8d74561DecodeGithubComNdRTechDbForumModels3(in *jlexer.Lexer, out *Forum) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "posts":
			if in.IsNull() {
				in.Skip()
				out.Posts = nil
			} else {
				if out.Posts == nil {
					out.Posts = new(int)
				}
				*out.Posts = int(in.Int())
			}
		case "slug":
			out.Slug = string(in.String())
		case "threads":
			if in.IsNull() {
				in.Skip()
				out.Threads = nil
			} else {
				if out.Threads == nil {
					out.Threads = new(int)
				}
				*out.Threads = int(in.Int())
			}
		case "title":
			out.Title = string(in.String())
		case "user":
			out.Moderator = string(in.String())
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
func easyjsonC8d74561EncodeGithubComNdRTechDbForumModels3(out *jwriter.Writer, in Forum) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"posts\":")
	if in.Posts == nil {
		out.RawString("null")
	} else {
		out.Int(int(*in.Posts))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"slug\":")
	out.String(string(in.Slug))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"threads\":")
	if in.Threads == nil {
		out.RawString("null")
	} else {
		out.Int(int(*in.Threads))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"title\":")
	out.String(string(in.Title))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"user\":")
	out.String(string(in.Moderator))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Forum) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8d74561EncodeGithubComNdRTechDbForumModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Forum) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8d74561EncodeGithubComNdRTechDbForumModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Forum) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8d74561DecodeGithubComNdRTechDbForumModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Forum) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8d74561DecodeGithubComNdRTechDbForumModels3(l, v)
}