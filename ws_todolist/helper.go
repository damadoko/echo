package main

func createID(c ClientResponse) int {
	// Create newID = lasted Todos ID + 1
	if index := len(c.Todos); index > 0 {
		return c.Todos[index].ID + 1 
	}
	return 0
}

func clearTodo(c ClientResponse) []Todo {
	j := 0
			for _, t := range c.Todos {
				if (t.Completed) {
					c.Todos[j] = t	
					j++
				}
			}
			return c.Todos[:j]
}