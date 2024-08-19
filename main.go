package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"

	"github.com/dghubble/ipnets"
)

func main() {
	code := run()
	if code != 0 {
		os.Exit(code)
	}
}

func run() int {
	var netmaskBits int
	flag.IntVar(&netmaskBits, "mask", 24, "bitmask used when generating subnet")
	flag.Parse()

	if netmaskBits < 1 || netmaskBits > 32 {
		fmt.Fprintln(os.Stderr, "ERROR: requested netmask size must be in the interval [1,32]")
		return 1
	}

	ranges := []string{
		// 172.17-31.x.x/16
		"172.17.0.0/16",
		"172.18.0.0/16",
		"172.19.0.0/16",
		"172.20.0.0/14",
		"172.24.0.0/14",
		"172.28.0.0/14",
		// 192.168.x.x/20
		"192.168.0.0/16",
		// 10.x.x.x/24
		"10.0.0.0/8",
	}
	shuffle(ranges)

	for _, rawCIDR := range ranges {
		_, cidr, err := net.ParseCIDR(rawCIDR)
		if err != nil || cidr.IP.To4() == nil {
			fmt.Fprintf(os.Stderr, "WARN: invalid IPv4 cidr string %q\n", cidr)
			continue
		}

		fmt.Println("mask", netmaskBits)
		maskSize, _ := cidr.Mask.Size()
		if netmaskBits < maskSize {
			continue // not big enough
		}
		shiftCount := netmaskBits - maskSize

		split, err := ipnets.SubnetShift(cidr, shiftCount)
		if err == nil && len(split) > 0 {
			shuffle(split)
			fmt.Println(split[0])
			return 0
		}
	}

	fmt.Fprintln(os.Stderr, "ERROR: could not get a free ipv4 private network")
	return 1
}

func shuffle[V any](list []V) {
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
}
