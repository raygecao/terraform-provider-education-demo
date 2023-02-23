package mockplatform

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"terraform-provider-education-demo/mockplatform/model"

	"gopkg.in/yaml.v3"
)

type Platform struct {
	storePath string
	mu        sync.RWMutex
}

var salaries = map[string]int{
	"math":    20000,
	"physic":  25000,
	"english": 15000,
}

type DataSet struct {
	Teachers []*model.Teacher `yaml:"teachers,omitempty"`
}

const path = "/tmp/platform-store.yaml"

func NewPlatform(user, passwd string) (*Platform, error) {
	return &Platform{storePath: fmt.Sprintf("/tmp/%s-%s-store.yaml", user, passwd)}, nil
}

func (p *Platform) getDataSet() *DataSet {
	p.mu.RLock()
	defer p.mu.RUnlock()
	f, err := os.Open(p.storePath)
	if err != nil {
		return &DataSet{}
	}
	var ds DataSet
	if err := yaml.NewDecoder(f).Decode(&ds); err != nil {
		panic(err)
	}
	return &ds
}

func (p *Platform) writeDataSet(ds *DataSet) {
	p.mu.Lock()
	defer p.mu.Unlock()
	f, _ := os.Create(p.storePath)
	if err := yaml.NewEncoder(f).Encode(ds); err != nil {
		panic(err)
	}
}

func (p *Platform) CreateTeacher(teacher *model.Teacher) (*model.Teacher, error) {
	ds := p.getDataSet()
	for _, t := range ds.Teachers {
		if t.ID == teacher.ID {
			return nil, fmt.Errorf("teacher with id %s exists", teacher.ID)
		}
	}
	if teacher.ID == 0 || teacher.Name == "" || teacher.Subject == "" {
		return nil, fmt.Errorf("miss required field")
	}
	salary, ok := salaries[teacher.Subject]
	if !ok {
		return nil, fmt.Errorf("unsupported subject %s", teacher.Subject)
	}
	if teacher.Salary == 0 {
		teacher.Salary = salary
	}

	ds.Teachers = append(ds.Teachers, teacher)
	sort.Slice(ds.Teachers, func(i, j int) bool {
		return ds.Teachers[i].ID < ds.Teachers[j].ID
	})
	p.writeDataSet(ds)
	return teacher, nil
}

func (p *Platform) GetTeacher(id int) (*model.Teacher, error) {
	ds := p.getDataSet()
	for _, t := range ds.Teachers {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, fmt.Errorf("teacher with id %s doesn't exist", id)
}

func (p *Platform) UpdateTeacher(teacher *model.Teacher) (*model.Teacher, error) {
	ds := p.getDataSet()
	t, err := p.GetTeacher(teacher.ID)
	if err != nil {
		return nil, err
	}
	if teacher.Subject != "" {
		salary, ok := salaries[teacher.Subject]
		if !ok {
			return nil, fmt.Errorf("unsupport subject %s", teacher.Subject)
		}
		t.Subject = teacher.Subject
		t.Salary = salary
	}
	if teacher.Salary != 0 {
		t.Salary = teacher.Salary
	}
	if teacher.Name != "" {
		t.Name = teacher.Name
	}
	for i := range ds.Teachers {
		if ds.Teachers[i].ID == teacher.ID {
			ds.Teachers[i] = teacher
		}
	}
	p.writeDataSet(ds)
	return t, nil
}

func (p *Platform) DeleteTeacher(id int) error {
	ds := p.getDataSet()
	if _, err := p.GetTeacher(id); err != nil {
		return err
	}
	for i, t := range ds.Teachers {
		if t.ID == id {
			ds.Teachers = append(ds.Teachers[:i], ds.Teachers[i+1:]...)
		}
	}
	p.writeDataSet(ds)
	return nil
}
