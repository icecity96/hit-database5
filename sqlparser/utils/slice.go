package utils

// InSlice checks given string in string slice or not.
func InSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

// InSliceIface checks given interface in interface slice.
func InSliceIface(v interface{}, sl []interface{}) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

// SliceIntersect returns slice that are present in all the slice1 and slice2.
func SliceIntersect(slice1, slice2 []string) (diffslice []string) {
	for _, v := range slice1 {
		if InSlice(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

// SliceDiff returns diff slice of slice1 - slice2.
func SliceDiff(slice1, slice2 []string) (diffslice []string) {
	for _, v := range slice1 {
		if !InSlice(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

// SliceMerge merges interface slices to one slice.
func SliceMerge(slice1, slice2 []string) (c []string) {
	c = append(slice1, slice2...)
	return
}

// SliceUnique cleans repeated values in slice.
func SliceUnique(slice []string) (uniqueslice []string) {
	for _, v := range slice {
		if !InSlice(v, uniqueslice) {
			uniqueslice = append(uniqueslice, v)
		}
	}
	return
}