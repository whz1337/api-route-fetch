package main

import (
	"reflect"
	"testing"
)

func TestSortDataIfTheDataIsAlreadySortedProperly(t *testing.T) {
	testData := correctOutput{}
	route := route{}
	testData.Source = "foo"
	route.Destination = "bar"
	route.Distance = 1.0
	route.Duration = 1.0

	testData.Routes = append(testData.Routes, route)

	route.Distance = 99.0
	route.Duration = 99.0

	testData.Routes = append(testData.Routes, route)

	ans := testData.Routes

	testData.Routes = sortData(testData.Routes)

	if(!reflect.DeepEqual(ans, testData.Routes)){
		t.Errorf("TestSortDataIfTheDataIsAlreadySortedProperly")
	}
}

func TestIfSortDataWorksProperlyOnOneRouteObjectOnly(t *testing.T){
	testData := correctOutput{}
	route := route{}
	testData.Source = "foo"
	route.Destination = "bar"
	route.Distance = 1.0
	route.Duration = 1.0

	testData.Routes = append(testData.Routes, route)

	ans := testData.Routes

	testData.Routes = sortData(testData.Routes)

	if(!reflect.DeepEqual(ans, testData.Routes)){
		t.Errorf("TestIfSortDataWorksProperlyOnOneRouteObjectOnly")
	}
}

func TestSortDataIfDurationIsTheSameButDistanceIsInWrongOrder(t *testing.T){
	testData := correctOutput{}
	ans := correctOutput{}
	route := route{}
	testData.Source = "foo"
	ans.Source = "foo"
	route.Destination = "bar"
	route.Distance = 2.0
	route.Duration = 1.0

	testData.Routes = append(testData.Routes, route)

	route.Distance = 1.0
	route.Duration = 1.0

	testData.Routes = append(testData.Routes, route)
	ans.Routes = append(ans.Routes, route)

	route.Distance = 2.0
	route.Duration = 1.0

	ans.Routes = append(ans.Routes, route)

	testData.Routes = sortData(testData.Routes)

	if(!reflect.DeepEqual(ans, testData)){
		t.Errorf("TestSortDataIfDurationIsTheSameButDistanceIsInWrongOrder")
	}

}

func TestSortDataIfThereAreThreeEntriesInWrongOrder(t *testing.T){
	testData := correctOutput{}
	ans := correctOutput{}
	route := route{}
	testData.Source = "foo"
	ans.Source = "foo"
	route.Destination = "bar"

	route.Distance = 1.0
	route.Duration = 2.0

	testData.Routes = append(testData.Routes, route)

	route.Distance = 1.0
	route.Duration = 1.0

	testData.Routes = append(testData.Routes, route)
	ans.Routes = append(ans.Routes, route)

	route.Distance = 1.0
	route.Duration = 3.0

	testData.Routes = append(testData.Routes, route)

	route.Distance = 1.0
	route.Duration = 2.0

	ans.Routes = append(ans.Routes, route)

	route.Distance = 1.0
	route.Duration = 3.0

	ans.Routes = append(ans.Routes, route)

	testData.Routes = sortData(testData.Routes)

	if(!reflect.DeepEqual(ans, testData)){
		t.Errorf("TestSortDataIfDurationIsTheSameButDistanceIsInWrongOrder")
	}
}

