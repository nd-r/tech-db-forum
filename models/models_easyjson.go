// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels(in *jlexer.Lexer, out *VoteDB) {
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
		case "ID":
			out.ID = int(in.Int())
		case "Nickname":
			out.Nickname = string(in.String())
		case "Thread_id":
			out.Thread_id = int(in.Int())
		case "Voice":
			out.Voice = int(in.Int())
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels(out *jwriter.Writer, in VoteDB) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"ID\":")
	out.Int(int(in.ID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Nickname\":")
	out.String(string(in.Nickname))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Thread_id\":")
	out.Int(int(in.Thread_id))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Voice\":")
	out.Int(int(in.Voice))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v VoteDB) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v VoteDB) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *VoteDB) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *VoteDB) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels1(in *jlexer.Lexer, out *Vote) {
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
		case "nickname":
			out.Nickname = string(in.String())
		case "voice":
			out.Voice = int(in.Int())
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels1(out *jwriter.Writer, in Vote) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"nickname\":")
	out.String(string(in.Nickname))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"voice\":")
	out.Int(int(in.Voice))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Vote) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Vote) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Vote) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Vote) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels1(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels2(in *jlexer.Lexer, out *UsersArr) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(UsersArr, 0, 8)
			} else {
				*out = UsersArr{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 *User
			if in.IsNull() {
				in.Skip()
				v1 = nil
			} else {
				if v1 == nil {
					v1 = new(User)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*v1).UnmarshalJSON(data))
				}
			}
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels2(out *jwriter.Writer, in UsersArr) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			if v3 == nil {
				out.RawString("null")
			} else {
				out.Raw((*v3).MarshalJSON())
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v UsersArr) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UsersArr) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UsersArr) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UsersArr) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels2(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels3(in *jlexer.Lexer, out *UserUpd) {
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
		case "about":
			if in.IsNull() {
				in.Skip()
				out.About = nil
			} else {
				if out.About == nil {
					out.About = new(string)
				}
				*out.About = string(in.String())
			}
		case "email":
			if in.IsNull() {
				in.Skip()
				out.Email = nil
			} else {
				if out.Email == nil {
					out.Email = new(string)
				}
				*out.Email = string(in.String())
			}
		case "fullname":
			if in.IsNull() {
				in.Skip()
				out.Fullname = nil
			} else {
				if out.Fullname == nil {
					out.Fullname = new(string)
				}
				*out.Fullname = string(in.String())
			}
		case "nickname":
			if in.IsNull() {
				in.Skip()
				out.Nickname = nil
			} else {
				if out.Nickname == nil {
					out.Nickname = new(string)
				}
				*out.Nickname = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels3(out *jwriter.Writer, in UserUpd) {
	out.RawByte('{')
	first := true
	_ = first
	if in.About != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"about\":")
		if in.About == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.About))
		}
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"email\":")
	if in.Email == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Email))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"fullname\":")
	if in.Fullname == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Fullname))
	}
	if in.Nickname != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"nickname\":")
		if in.Nickname == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Nickname))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserUpd) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserUpd) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserUpd) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserUpd) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels3(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels4(in *jlexer.Lexer, out *User) {
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
		case "about":
			out.About = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "fullname":
			out.Fullname = string(in.String())
		case "nickname":
			out.Nickname = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels4(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	if in.About != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"about\":")
		out.String(string(in.About))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"email\":")
	out.String(string(in.Email))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"fullname\":")
	out.String(string(in.Fullname))
	if in.Nickname != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"nickname\":")
		out.String(string(in.Nickname))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels4(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels5(in *jlexer.Lexer, out *TreadArr) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(TreadArr, 0, 8)
			} else {
				*out = TreadArr{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v4 *Thread
			if in.IsNull() {
				in.Skip()
				v4 = nil
			} else {
				if v4 == nil {
					v4 = new(Thread)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*v4).UnmarshalJSON(data))
				}
			}
			*out = append(*out, v4)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels5(out *jwriter.Writer, in TreadArr) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v5, v6 := range in {
			if v5 > 0 {
				out.RawByte(',')
			}
			if v6 == nil {
				out.RawString("null")
			} else {
				out.Raw((*v6).MarshalJSON())
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v TreadArr) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TreadArr) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TreadArr) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TreadArr) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels5(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels6(in *jlexer.Lexer, out *ThreadUpdate) {
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
			if in.IsNull() {
				in.Skip()
				out.Message = nil
			} else {
				if out.Message == nil {
					out.Message = new(string)
				}
				*out.Message = string(in.String())
			}
		case "title":
			if in.IsNull() {
				in.Skip()
				out.Title = nil
			} else {
				if out.Title == nil {
					out.Title = new(string)
				}
				*out.Title = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels6(out *jwriter.Writer, in ThreadUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"message\":")
	if in.Message == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Message))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"title\":")
	if in.Title == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Title))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels6(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels7(in *jlexer.Lexer, out *Thread) {
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
			out.Forum_slug = string(in.String())
		case "author":
			out.User_nick = string(in.String())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels7(out *jwriter.Writer, in Thread) {
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
	out.String(string(in.Forum_slug))
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
	out.Raw((in.Created).MarshalJSON())
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
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Thread) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Thread) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Thread) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels7(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels8(in *jlexer.Lexer, out *Status) {
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
		case "forum":
			out.Forum = int(in.Int())
		case "post":
			out.Post = int(in.Int())
		case "thread":
			out.Thread = int(in.Int())
		case "user":
			out.User = int(in.Int())
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels8(out *jwriter.Writer, in Status) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"forum\":")
	out.Int(int(in.Forum))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"post\":")
	out.Int(int(in.Post))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"thread\":")
	out.Int(int(in.Thread))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"user\":")
	out.Int(int(in.User))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Status) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Status) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Status) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Status) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels8(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels9(in *jlexer.Lexer, out *PostUpdate) {
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
			if in.IsNull() {
				in.Skip()
				out.Message = nil
			} else {
				if out.Message == nil {
					out.Message = new(string)
				}
				*out.Message = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels9(out *jwriter.Writer, in PostUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"message\":")
	if in.Message == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels9(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels10(in *jlexer.Lexer, out *PostDetails) {
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
		case "author":
			if in.IsNull() {
				in.Skip()
				out.AuthorDetails = nil
			} else {
				if out.AuthorDetails == nil {
					out.AuthorDetails = new(User)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.AuthorDetails).UnmarshalJSON(data))
				}
			}
		case "forum":
			if in.IsNull() {
				in.Skip()
				out.ForumDetails = nil
			} else {
				if out.ForumDetails == nil {
					out.ForumDetails = new(Forum)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.ForumDetails).UnmarshalJSON(data))
				}
			}
		case "post":
			if in.IsNull() {
				in.Skip()
				out.PostDetails = nil
			} else {
				if out.PostDetails == nil {
					out.PostDetails = new(Post)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.PostDetails).UnmarshalJSON(data))
				}
			}
		case "thread":
			if in.IsNull() {
				in.Skip()
				out.ThreadDetails = nil
			} else {
				if out.ThreadDetails == nil {
					out.ThreadDetails = new(Thread)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.ThreadDetails).UnmarshalJSON(data))
				}
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels10(out *jwriter.Writer, in PostDetails) {
	out.RawByte('{')
	first := true
	_ = first
	if in.AuthorDetails != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"author\":")
		if in.AuthorDetails == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.AuthorDetails).MarshalJSON())
		}
	}
	if in.ForumDetails != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"forum\":")
		if in.ForumDetails == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.ForumDetails).MarshalJSON())
		}
	}
	if in.PostDetails != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"post\":")
		if in.PostDetails == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.PostDetails).MarshalJSON())
		}
	}
	if in.ThreadDetails != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"thread\":")
		if in.ThreadDetails == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.ThreadDetails).MarshalJSON())
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostDetails) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostDetails) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostDetails) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostDetails) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels10(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels11(in *jlexer.Lexer, out *PostArr) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(PostArr, 0, 8)
			} else {
				*out = PostArr{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v7 *Post
			if in.IsNull() {
				in.Skip()
				v7 = nil
			} else {
				if v7 == nil {
					v7 = new(Post)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*v7).UnmarshalJSON(data))
				}
			}
			*out = append(*out, v7)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels11(out *jwriter.Writer, in PostArr) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in {
			if v8 > 0 {
				out.RawByte(',')
			}
			if v9 == nil {
				out.RawString("null")
			} else {
				out.Raw((*v9).MarshalJSON())
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v PostArr) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels11(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostArr) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels11(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostArr) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels11(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostArr) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels11(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels12(in *jlexer.Lexer, out *Post) {
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
			out.Id = int(in.Int())
		case "author":
			out.User_nick = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "created":
			if in.IsNull() {
				in.Skip()
				out.Created = nil
			} else {
				if out.Created == nil {
					out.Created = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.Created).UnmarshalJSON(data))
				}
			}
		case "forum":
			out.Forum_slug = string(in.String())
		case "thread":
			out.Thread_id = int(in.Int())
		case "isEdited":
			out.Is_edited = bool(in.Bool())
		case "parent":
			out.Parent = int64(in.Int64())
		case "Parents":
			if in.IsNull() {
				in.Skip()
				out.Parents = nil
			} else {
				in.Delim('[')
				if out.Parents == nil {
					if !in.IsDelim(']') {
						out.Parents = make([]int64, 0, 8)
					} else {
						out.Parents = []int64{}
					}
				} else {
					out.Parents = (out.Parents)[:0]
				}
				for !in.IsDelim(']') {
					var v10 int64
					v10 = int64(in.Int64())
					out.Parents = append(out.Parents, v10)
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels12(out *jwriter.Writer, in Post) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	out.Int(int(in.Id))
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
	out.RawString("\"message\":")
	out.String(string(in.Message))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"created\":")
	if in.Created == nil {
		out.RawString("null")
	} else {
		out.Raw((*in.Created).MarshalJSON())
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"forum\":")
	out.String(string(in.Forum_slug))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"thread\":")
	out.Int(int(in.Thread_id))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"isEdited\":")
	out.Bool(bool(in.Is_edited))
	if in.Parent != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"parent\":")
		out.Int64(int64(in.Parent))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Parents\":")
	if in.Parents == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v11, v12 := range in.Parents {
			if v11 > 0 {
				out.RawByte(',')
			}
			out.Int64(int64(v12))
		}
		out.RawByte(']')
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Post) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels12(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Post) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels12(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Post) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels12(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Post) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels12(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels13(in *jlexer.Lexer, out *Forum) {
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
		case "slug":
			out.Slug = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "user":
			out.Moderator = string(in.String())
		case "threads":
			out.Threads = int(in.Int())
		case "posts":
			out.Posts = int(in.Int())
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels13(out *jwriter.Writer, in Forum) {
	out.RawByte('{')
	first := true
	_ = first
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
	out.RawString("\"title\":")
	out.String(string(in.Title))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"user\":")
	out.String(string(in.Moderator))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"threads\":")
	out.Int(int(in.Threads))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"posts\":")
	out.Int(int(in.Posts))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Forum) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels13(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Forum) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels13(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Forum) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels13(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Forum) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels13(l, v)
}
func easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels14(in *jlexer.Lexer, out *ErrorStr) {
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
func easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels14(out *jwriter.Writer, in ErrorStr) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"message\":")
	out.String(string(in.Message))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ErrorStr) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels14(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ErrorStr) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComNdRTechDbForumModels14(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ErrorStr) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels14(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ErrorStr) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComNdRTechDbForumModels14(l, v)
}
