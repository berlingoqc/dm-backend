package tr_test

import "testing"

func TestPipeline(t *testing.T) {
	defer destroyEven()

	prepareEven()

	initModule(t)
	createPipeline(t)
	getPipeline(t)

}

func prepareEven() {

}

func destroyEven() {

}

func initModule(t *testing.T) {

}

func createPipeline(t *testing.T) {

}

func getPipeline(t *testing.T) {

}
