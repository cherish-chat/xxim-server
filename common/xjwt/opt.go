package xjwt

type Option struct {
	Platform    *string
	DeviceId    *string
	DeviceModel *string
	Data        *string
}
type OptionFunc func(o *Option)

func WithPlatform(platform string) OptionFunc {
	return func(o *Option) {
		o.Platform = &platform
	}
}
func WithDeviceId(deviceId string) OptionFunc {
	return func(o *Option) {
		o.DeviceId = &deviceId
	}
}
func WithDeviceModel(deviceModel string) OptionFunc {
	return func(o *Option) {
		o.DeviceModel = &deviceModel
	}
}
func WithData(data string) OptionFunc {
	return func(o *Option) {
		o.Data = &data
	}
}
