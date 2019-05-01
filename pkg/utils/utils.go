package utils

import (
	"fmt"
	"net"
	"reflect"
	"sort"

	"github.com/Microsoft/KubeDevice-API/pkg/kdlog"
	"github.com/Microsoft/KubeDevice-API/pkg/types"
)

// Logb tells whether to log
var Logb func(int) bool

// Logf provides logging functionality inside plugins
var Logf func(int, string, ...interface{})

// Warningf provides logging functionality inside plugins
var Warningf func(string, ...interface{})

// Errorf provides logginf functionality inside plugins
var Errorf func(string, ...interface{})

func LogV(level int) bool {
	return bool(kdlog.V(kdlog.Level(level)))
}

func Log(level int, format string, args ...interface{}) {
	if kdlog.V(kdlog.Level(level)) {
		str := fmt.Sprintf(format, args...)
		kdlog.InfoDepth(1, str)
	}
}

func Error(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	kdlog.ErrorDepth(1, str)
}

func Warning(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	kdlog.WarningDepth(1, str)
}

func init() {
	Logb = LogV
	Logf = Log
	Errorf = Error
	Warningf = Warning
}

func LocalIPsWithoutLoopback() ([]net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("could not list network interfaces: %v", err)
	}
	var ips []net.IP
	for _, i := range interfaces {
		addresses, err := i.Addrs()
		if err != nil {
			return nil, fmt.Errorf("could not list the addresses for network interface %v: %v\n", i, err)
		}
		for _, address := range addresses {
			switch v := address.(type) {
			case *net.IPNet:
				if !v.IP.IsLoopback() {
					ips = append(ips, v.IP)
				}
			}
		}
	}
	return ips, nil
}

// sorted string keys
func SortedStringKeys(x interface{}) []string {
	t := reflect.TypeOf(x)
	keys := []string{}
	if t.Kind() == reflect.Map {
		mv := reflect.ValueOf(x)
		keysV := mv.MapKeys()
		for _, val := range keysV {
			keys = append(keys, val.String())
		}
		sort.Strings(keys)
		return keys
	}
	panic("Not a map")
}

func CompareContainer(cont0 *types.ContainerInfo, cont1 *types.ContainerInfo) bool {
	if true {
		if !reflect.DeepEqual(cont0.KubeRequests, cont1.KubeRequests) {
			fmt.Printf("KubeReqs don't match\n0:\n%v\n1:\n%v\n", cont0.KubeRequests, cont1.KubeRequests)
			return false
		}
		if !reflect.DeepEqual(cont0.Requests, cont1.Requests) {
			fmt.Printf("Reqs don't match\n0:\n%v\n1:\n%v\n", cont0.Requests, cont1.Requests)
			return false
		}
		if !reflect.DeepEqual(cont0.DevRequests, cont1.DevRequests) {
			fmt.Printf("DevReqs don't match\n0:\n%v\n1:\n%v\n", cont0.DevRequests, cont1.DevRequests)
			return false
		}
		if !reflect.DeepEqual(cont0.AllocateFrom, cont1.AllocateFrom) {
			fmt.Printf("AllocateFrom don't match\n0:\n%v\n1:\n%v\n", cont0.AllocateFrom, cont1.AllocateFrom)
			return false
		}
		if !reflect.DeepEqual(cont0.Scorer, cont1.Scorer) {
			fmt.Printf("Scorer don't match\n0:\n%v\n1:\n%v\n", cont0.Scorer, cont1.Scorer)
			return false
		}
	}
	return true
}

func CompareContainers(conts0 map[string]types.ContainerInfo, conts1 map[string]types.ContainerInfo) bool {
	allSame := true
	for contName0, cont0 := range conts0 {
		cont1, ok := conts1[contName0]
		if !ok {
			fmt.Printf("1 does not have container %s\n", contName0)
			allSame = false
		} else {
			fmt.Printf("Compare container %s\n", contName0)
			allSame = allSame && CompareContainer(&cont0, &cont1)
		}
	}
	for contName1, _ := range conts1 {
		_, ok := conts0[contName1]
		if !ok {
			fmt.Printf("0 does not have container %s\n", contName1)
			allSame = false
		}
	}
	return allSame
}

func ComparePod(pod0 *types.PodInfo, pod1 *types.PodInfo) bool {
	allSame := true
	if pod0.Name != pod1.Name {
		fmt.Printf("Name does not match %s %s\n", pod0.Name, pod1.Name)
		allSame = false
	}
	if pod0.NodeName != pod1.NodeName {
		fmt.Printf("Nodename does not match %s %s\n", pod0.NodeName, pod1.NodeName)
		allSame = false
	}
	if !reflect.DeepEqual(pod0.Requests, pod1.Requests) {
		fmt.Printf("PodRequests don't match\n0:\n%+v\n1:\n%v\n", pod0.Requests, pod1.Requests)
		allSame = false
	}
	allSame = allSame && CompareContainers(pod0.InitContainers, pod1.InitContainers)
	allSame = allSame && CompareContainers(pod0.RunningContainers, pod1.RunningContainers)
	return allSame
}

func CompareNode(node0 *types.NodeInfo, node1 *types.NodeInfo) bool {
	allSame := true

	if node0.Name != node1.Name {
		fmt.Printf("Name does not match %s %s", node0.Name, node1.Name)
		allSame = false
	}
	if !reflect.DeepEqual(node0.Capacity, node1.Capacity) {
		fmt.Printf("Capacity not same\n0\n:%+v\n1:\n:%v\n", node0.Capacity, node1.Capacity)
		allSame = false
	}
	if !reflect.DeepEqual(node0.Allocatable, node1.Allocatable) {
		fmt.Printf("Allocatable not same\n0\n:%+v\n1:\n:%v\n", node0.Allocatable, node1.Allocatable)
		allSame = false
	}
	if !reflect.DeepEqual(node0.Used, node1.Used) {
		fmt.Printf("Used not same\n0\n:%+v\n1:\n:%v\n", node0.Used, node1.Used)
		allSame = false
	}
	if !reflect.DeepEqual(node0.Scorer, node1.Scorer) {
		fmt.Printf("Scorer not same\n0\n:%+v\n1:\n:%v\n", node0.Scorer, node1.Scorer)
		allSame = false
	}
	if !reflect.DeepEqual(node0.KubeCap, node1.KubeCap) {
		fmt.Printf("KubeCap not same\n0\n:%+v\n1:\n:%v\n", node0.KubeCap, node1.KubeCap)
		allSame = false
	}
	if !reflect.DeepEqual(node0.KubeAlloc, node1.KubeAlloc) {
		fmt.Printf("KubeAlloc not same\n0\n:%+v\n1:\n:%v\n", node0.KubeAlloc, node1.KubeAlloc)
		allSame = false
	}

	return allSame
}
