package main

import "fmt"

// Resource represents a system resource with an ID and data
type Resource struct {
	ID   int
	Data string
}

// ResourceHolder contains a pointer to a Resource
type ResourceHolder struct {
	Name     string
	Resource *Resource // Pointer to Resource
}

func main() {
	// Create some resources
	resource1 := &Resource{ID: 1, Data: "Resource 1 Data"}
	resource2 := &Resource{ID: 2, Data: "Resource 2 Data"}

	// Create ResourceHolders that point to these resources
	holder1 := ResourceHolder{
		Name:     "Holder 1",
		Resource: resource1,
	}

	holder2 := ResourceHolder{
		Name:     "Holder 2", 
		Resource: resource2,
	}

	// Create a map using ResourceHolder as key
	holderStatus := make(map[ResourceHolder]string)

	// Add holders to the map
	holderStatus[holder1] = "Active"
	holderStatus[holder2] = "Standby"

	// Print initial map state
	fmt.Println("Initial holder status:")
	for holder, status := range holderStatus {
		fmt.Printf("%s (Resource ID: %d): %s\n", 
			holder.Name, 
			holder.Resource.ID, 
			status)
	}

	// Create a new holder pointing to the same resource
	holder3 := ResourceHolder{
		Name:     "Holder 3",
		Resource: resource1, // Points to same resource as holder1
	}

	// Add to map
	holderStatus[holder3] = "Active"

	fmt.Println("\nAfter adding holder3 (points to same resource as holder1):")
	for holder, status := range holderStatus {
		fmt.Printf("%s (Resource ID: %d): %s\n", 
			holder.Name, 
			holder.Resource.ID, 
			status)
	}

	// Modify the resource through one holder
	holder1.Resource.Data = "Modified Resource 1 Data"

	// The modification is visible through all holders pointing to the same resource
	fmt.Println("\nAfter modifying resource through holder1:")
	fmt.Printf("Through holder1: %s\n", holder1.Resource.Data)
	fmt.Printf("Through holder3: %s\n", holder3.Resource.Data)
}
