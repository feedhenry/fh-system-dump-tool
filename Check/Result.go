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

	projectData := map[string]map[string]string{}

	for _, item := range c.Info {
		if _, ok := projectData[item.Namespace]; ! ok {
			projectData[item.Namespace] = map[string]string{}
		}
		projectData[item.Namespace][item.ObjectName] = item.Entry
	}

	for projectName, project := range projectData {
		fmt.Println("	Project: " + projectName)
		for podName, msg := range project {
			fmt.Println("		Pod (" + podName + "): " + msg)
		}
	}
}
