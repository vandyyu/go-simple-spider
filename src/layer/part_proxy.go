package layer
type IPartProxy interface{
	IPart
	Forward() (IPartProxy, interface{})
	InitPartProxy()
}

