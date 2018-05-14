package hyperion

import marathon "github.com/gambol99/go-marathon"

type marathonApp struct {
	gomApp *marathon.Application
}

type marathonAppDeleteRequest struct {
	appID string
}

func (m *marathonApp) ID() string {
	return m.gomApp.ID
}

func (m *marathonApp) SetID(id string) *marathonApp {
	m.gomApp.ID = id
	return m
}

func (m *marathonApp) SetCount(count int) *marathonApp {
	m.gomApp.Count(count)
	return m
}

func (m *marathonApp) SetDockerImage(dockerImage string) *marathonApp {
	m.gomApp.Container.Docker.Container(dockerImage)
	return m
}

func (m *marathonApp) NewDeleteRequest() *marathonAppDeleteRequest {
	return &marathonAppDeleteRequest{appID: m.ID()}
}
