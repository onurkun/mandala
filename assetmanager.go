package mandala

import (
	"io"
	"unsafe"

	"git.tideland.biz/goas/loop"
)

var (
	// The path in which the framework will search for resources.
	AssetPath string = "android"
)

// The type of object returned by the AssetManager in response of a
// LoadAssetRequest request.
type LoadAssetResponse struct {
	// The buffer containing the resource. Please note that, in
	// case io.Reader is os.File, it's a client code
	// responsibility to close it. This usually happens in the
	// xorg version of your program. Maybe this API will change in
	// the future, though.
	Buffer io.Reader

	// The error eventually generated by the reading operation.
	Error error
}

// The type of request to send to AssetManager in order to receive a
// LoadAssetResponse.
type LoadAssetRequest struct {
	// The path of the resource file, for example
	// "res/drawable/gopher.png". AssetPath will be prefixed to
	// Filename.
	Filename string

	// Response is a channel from which receive
	Response chan LoadAssetResponse
}

// Run runs the nativeEventsLoop.
// The loop handles native input events.
func assetLoopFunc(activity chan unsafe.Pointer, request chan interface{}) loop.LoopFunc {
	var act unsafe.Pointer
	return func(l loop.Loop) error {
		for {
			select {
			case act = <-activity:
			case untypedRequest := <-request:
				switch request := untypedRequest.(type) {
				case LoadAssetRequest:
					file, err := loadAsset(act, request.Filename)
					request.Response <- LoadAssetResponse{file, err}
				}
			}
		}
	}
}

// func LoadAsset(filename string) <-chan io.Reader {
// 	command := LoadAssetCommand{
// 		Filename: filename,
// 		Buffer:   make(chan io.Reader),
// 	}
// 	Assets <- command
// 	return command.Buffer
// }
