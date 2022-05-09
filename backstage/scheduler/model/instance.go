package model

/*
	@model: StreamInstance
	@description:
		represent an instance of stream application
*/
type StreamInstance struct {
	*StreamApplication
	InstanceID string `json:"instance_id"`
}

/*
	@model: StreamInstanceRoom
	@description:
		room of a initialized stream application instance
*/
type StreamInstanceRoom struct {
	*StreamInstance

	Provider     *Provider
	ConsumerList map[string]*Consumer

	SelectedDepositary      *Depositary
	PotentialDepositaryList map[string]*Depositary

	SelectedFilestore      *Filestore
	PotentialFilestoreList map[string]*Filestore
}
