package layer
type IPart interface{
	GetName() string
	InitPart()
	Run() (IPart, interface{})
	SetPipeline(pipeline IPipeline)
	SetUnit(u IUnit)
	GetUnit() IUnit
}
