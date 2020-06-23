package main

func createID(c ClientResponse) int {
	// Create newID = lasted Todos ID + 1
	if index := len(c.Todos); index > 0 {
		return c.Todos[index].ID + 1 
	}
	return 0
}

func (c ClientResponse) deleteTodo(id int)  {
	for i, t := range c.Todos {
		if t.ID == id {
			c.Todos = append(c.Todos[:i], c.Todos[i+1:0]...)
		}
	}	
}

func (c ClientResponse) clearTodo() {
	j := 0
			for _, t := range c.Todos {
				if (t.Completed) {
					c.Todos[j] = t	
					j++
				}
			}
			c.Todos = c.Todos[:j]
}

func (c ClientResponse) toggleComplete(id int)  {
	for i, t := range c.Todos {
		if t.ID == id {
			c.Todos[i].Completed = !c.Todos[i].Completed
		}
	}	
}