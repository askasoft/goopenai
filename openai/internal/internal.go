package internal

import "github.com/askasoft/pango/doc/jsonx"

func ToString(o any) string {
	return jsonx.Prettify(o)
}
