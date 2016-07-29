package Check
import "fmt"


type CheckResult interface {
	Output()
}

type Result struct {
	Status    int
	CheckName string
	Info      []Info
}

func (c *Result) Output() {
	fmt.Println(c.CheckName + " results: ")
	if c.Status == 0 {
		fmt.Println("	✔ - Issue not detected.")
		return
	}

	for _, item := range c.Info {
		fmt.Println("	✘ - " + item.ObjectName + " has ImagePullBackOff issue")
	}
}
