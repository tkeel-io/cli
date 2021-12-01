package kubernetes

import (
	"context"
	"fmt"
)

func ExampleListPluginPods() {
	client, err := Client()
	if err != nil {
		fmt.Println(err)
	}
	{
		ret, err := ListPluginPods(client)
		if err != nil {
			panic(err)
		}
		fmt.Println("pods:", len(ret))

		ret, err = ListPluginPods(client, "rudder")
		if err != nil {
			panic(err)
		}
		fmt.Println("pods:", len(ret))

		for _, p := range ret {
			app := p.App()
			if app.AppID == "rudder" {
				res := app.Request(client.CoreV1().RESTClient().Get()).Suffix("v1/plugins")
				ret := res.Do(context.TODO())
				raw, err := ret.Raw()
				fmt.Println(string(raw), err)
			}
		}
	}

	// Output:
	// pods: 4
	// pods: 1
	// {"pluginList":[{"id":"keel-echo","tkeel_version":"v0.2.0","secret":"changeme","register_timestamp":1637909273}]} <nil>
}

func ExampleRegisterPlugins() {
	err := Register("keel-echo")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Unregister("keel-echo")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	//
}

func ExampleUnregisterPlugins() {
	err := Register("keel-echo")
	if err != nil {
		fmt.Println(err)
	}
	_, err = Unregister("keel-echo")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	//
}
