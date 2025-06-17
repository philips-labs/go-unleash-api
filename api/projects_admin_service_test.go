package api

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/philips-labs/go-unleash-api/mocks"
)

var (
	projectsService *ProjectsService
)

func init() {
	projectsService = &ProjectsService{
		client: &ApiClient{
			client:    &mocks.MockClient{},
			apiUrl:    &url.URL{Path: "local"},
			authToken: "myToken",
		},
	}

}

func TestProjectsService_GetProjectById(t *testing.T) {
	httpResponseMocks := make(map[string]*http.Response)
	httpResponseMocks["success"] = createHttpResponseMock(200, `{"name":"Default","description":"Default project","environments":[
        {
            "environment": "development"
        },
        {
            "environment": "production"
        }
    ],"features":[]}`, "GET")
	httpResponseMocks["notfound"] = createHttpResponseMock(404, `{"name":"NotFoundError"`, "GET")
	type args struct {
		projectId string
	}
	tests := []struct {
		name           string
		p              *ProjectsService
		args           args
		mockedResponse *http.Response
		wantProject    *ProjectDetails
		wantResponse   *Response
		wantErr        bool
	}{
		{
			"ReturnsProjectDetails",
			projectsService,
			args{
				projectId: "default",
			},
			httpResponseMocks["success"],
			&ProjectDetails{
				Name:         "Default",
				Description:  "Default project",
				Environments: []ProjEnvironment{{Environment: "development"}, {Environment: "production"}},
			},
			&Response{Response: httpResponseMocks["success"]},
			false,
		},
		{
			"ReturnsError",
			projectsService,
			args{
				projectId: "foo",
			},
			httpResponseMocks["notfound"],
			nil,
			&Response{Response: httpResponseMocks["notfound"]},
			true,
		},
	}
	for _, tt := range tests {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return tt.mockedResponse, nil
		}
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.p.GetProjectById(tt.args.projectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectsService.GetProjectById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantProject) {
				t.Errorf("ProjectsService.GetProjectById() got = %v, want %v", got, tt.wantProject)
			}
			if !reflect.DeepEqual(got1, tt.wantResponse) {
				t.Errorf("ProjectsService.GetProjectById() got1 = %v, want %v", got1, tt.wantResponse)
			}
		})
	}
}

func TestProjectsService_CreateProject(t *testing.T) {
	httpResponseMocks := make(map[string]*http.Response)
	httpResponseMocks["success"] = createHttpResponseMock(200, `{"id":"Default","name":"Default project","description":"Default project"}`, "POST")
	httpResponseMocks["notfound"] = createHttpResponseMock(404, `{"name":"NotFoundError"`, "POST")
	type args struct {
		Id          string
		Name        string
		Description string
	}
	tests := []struct {
		name           string
		p              *ProjectsService
		args           args
		mockedResponse *http.Response
		wantProject    *CreateProjectResponse
		wantResponse   *Response
		wantErr        bool
	}{
		{
			"SuccessfulRequest",
			projectsService,
			args{
				Id:          "Default",
				Name:        "Default",
				Description: "Default project",
			},
			httpResponseMocks["success"],
			&CreateProjectResponse{
				Id:          "Default",
				Name:        "Default project",
				Description: "Default project",
			},
			&Response{Response: httpResponseMocks["success"]},
			false,
		},
		{
			"ReturnsError",
			projectsService,
			args{
				Id:          "Default",
				Name:        "Default",
				Description: "Default project",
			},
			httpResponseMocks["notfound"],
			nil,
			&Response{Response: httpResponseMocks["notfound"]},
			true,
		},
	}
	for _, tt := range tests {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return tt.mockedResponse, nil
		}
		t.Run(tt.name, func(t *testing.T) {
			project := Project{
				Id:          tt.args.Id,
				Name:        tt.args.Name,
				Description: tt.args.Description,
			}
			got, got1, err := tt.p.CreateProject(project)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectsService.CreateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantProject) {
				t.Errorf("ProjectsService.CreateProject() got = %v, want %v", got, tt.wantProject)
			}
			if !reflect.DeepEqual(got1, tt.wantResponse) {
				t.Errorf("ProjectsService.CreateProject() got1 = %v, want %v", got1, tt.wantResponse)
			}
		})
	}
}

func TestProjectsService_UpdateProject(t *testing.T) {
	httpResponseMocks := make(map[string]*http.Response)
	httpResponseMocks["success"] = createHttpResponseMock(200, `{"id":"Default","name":"Default project","description":"Default project"}`, "PUT")
	httpResponseMocks["notfound"] = createHttpResponseMock(404, `{"name":"NotFoundError"`, "PUT")

	type args struct {
		Id          string
		Name        string
		Description string
	}

	tests := []struct {
		name           string
		p              *ProjectsService
		args           args
		mockedResponse *http.Response
		wantProject    *CreateProjectResponse
		wantResponse   *Response
		wantErr        bool
	}{
		{
			"SuccessfulRequest",
			projectsService,
			args{
				Id:          "Default",
				Name:        "Default",
				Description: "Default project",
			},
			httpResponseMocks["success"],
			&CreateProjectResponse{
				Id:          "Default",
				Name:        "Default project",
				Description: "Default project",
			},
			&Response{Response: httpResponseMocks["success"]},
			false,
		},
		{
			"ReturnsError",
			projectsService,
			args{
				Id:          "Default",
				Name:        "Default",
				Description: "Default project",
			},
			httpResponseMocks["notfound"],
			nil,
			&Response{Response: httpResponseMocks["notfound"]},
			true,
		},
	}
	for _, tt := range tests {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return tt.mockedResponse, nil
		}
		t.Run(tt.name, func(t *testing.T) {
			project := Project{
				Id:          tt.args.Id,
				Name:        tt.args.Name,
				Description: tt.args.Description,
			}
			got, got1, err := tt.p.UpdateProject(tt.args.Id, project)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectsService.CreateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantProject) {
				t.Errorf("ProjectsService.CreateProject() got = %v, want %v", got, tt.wantProject)
			}
			if !reflect.DeepEqual(got1, tt.wantResponse) {
				t.Errorf("ProjectsService.CreateProject() got1 = %v, want %v", got1, tt.wantResponse)
			}
		})
	}
}

func TestProjectsService_DeleteProject(t *testing.T) {
	httpResponseMocks := make(map[string]*http.Response)
	httpResponseMocks["success"] = createHttpResponseMock(204, "", "DELETE")
	httpResponseMocks["notfound"] = createHttpResponseMock(404, `{"name":"NotFoundError"}`, "DELETE")

	type args struct {
		projectId string
	}

	tests := []struct {
		name           string
		p              *ProjectsService
		args           args
		mockedResponse *http.Response
		wantResponse   *Response
		wantErr        bool
	}{
		{
			"SuccessfulRequest",
			projectsService,
			args{"Default"},
			httpResponseMocks["success"],
			&Response{Response: httpResponseMocks["success"]},
			false,
		},
		{
			"ReturnsError",
			projectsService,
			args{"InvalidProjectID"},
			httpResponseMocks["notfound"],
			&Response{Response: httpResponseMocks["notfound"]},
			true,
		},
	}

	for _, tt := range tests {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return tt.mockedResponse, nil
		}

		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.p.DeleteProject(tt.args.projectId)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectsService.DeleteProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestProjectsService_AddUserProject(t *testing.T) {
	httpResponseMocks := make(map[string]*http.Response)
	httpResponseMocks["success"] = createHttpResponseMock(200, `{"status":"success"}`, "POST")
	httpResponseMocks["badrequest"] = createHttpResponseMock(404, `{"status":"BadRequestError"}`, "POST")

	type args struct {
		userId    int
		projectId string
		roleId    int
	}
	tests := []struct {
		name           string
		p              *ProjectsService
		args           args
		mockedResponse *http.Response
		wantProject    *AddUserRoleResponse
		wantResponse   *Response
		wantErr        bool
	}{
		{
			"SuccessfulRequest",
			projectsService,
			args{
				userId:    100,
				projectId: "default",
				roleId:    3,
			},
			httpResponseMocks["success"],
			&AddUserRoleResponse{
				UserId:    100,
				ProjectId: "default",
				RoleId:    3,
			},
			&Response{Response: httpResponseMocks["success"]},
			false,
		},
		{
			"ReturnsError",
			projectsService,
			args{
				userId:    100,
				projectId: "default",
				roleId:    3,
			},
			httpResponseMocks["badrequest"],
			nil,
			&Response{Response: httpResponseMocks["badrequest"]},
			true,
		},
	}
	for _, tt := range tests {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return tt.mockedResponse, nil
		}
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := tt.p.AddUserProject(tt.args.userId, tt.args.projectId, tt.args.roleId)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectsService.AddUserRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// no additional checks for now as the response has no body.
			// ref: https://docs.getunleash.io/reference/api/legacy/unleash/admin/features-v2#add-a-user-to-a-project
		})
	}
}

func TestProjectsService_UpdateUserProject(t *testing.T) {
	httpResponseMocks := make(map[string]*http.Response)
	httpResponseMocks["success"] = createHttpResponseMock(200, `{"status":"success"}`, "POST")
	httpResponseMocks["badrequest"] = createHttpResponseMock(404, `{"status":"BadRequestError"}`, "POST")

	type args struct {
		UserId    int
		ProjectId string
		RoleId    int
	}
	tests := []struct {
		name           string
		p              *ProjectsService
		args           args
		mockedResponse *http.Response
		wantProject    *AddUserRoleResponse
		wantResponse   *Response
		wantErr        bool
	}{
		{
			"SuccessfulRequest",
			projectsService,
			args{
				UserId:    100,
				ProjectId: "default",
				RoleId:    3,
			},
			httpResponseMocks["success"],
			&AddUserRoleResponse{
				UserId:    100,
				ProjectId: "default",
				RoleId:    3,
			},
			&Response{Response: httpResponseMocks["success"]},
			false,
		},
		{
			"ReturnsError",
			projectsService,
			args{
				UserId:    100,
				ProjectId: "default",
				RoleId:    3,
			},
			httpResponseMocks["badrequest"],
			nil,
			&Response{Response: httpResponseMocks["badrequest"]},
			true,
		},
	}
	for _, tt := range tests {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return tt.mockedResponse, nil
		}
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := tt.p.UpdateUserProject(tt.args.ProjectId, tt.args.UserId, tt.args.RoleId)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectsService.UpdateUserProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// no additional checks for now as the response has no body.
			// ref: https://docs.getunleash.io/reference/api/legacy/unleash/admin/features-v2#change-a-users-role-in-a-project
		})
	}
}

func TestProjectsService_DeleteUserProject(t *testing.T) {
	httpResponseMocks := make(map[string]*http.Response)
	httpResponseMocks["success"] = createHttpResponseMock(200, `{"status":"success"}`, "POST")
	httpResponseMocks["badrequest"] = createHttpResponseMock(400, `{"status":"BadRequestError"}`, "POST")

	type args struct {
		UserId    int
		ProjectId string
		RoleId    int
	}
	tests := []struct {
		name           string
		p              *ProjectsService
		args           args
		mockedResponse *http.Response
		wantProject    *AddUserRoleResponse
		wantResponse   *Response
		wantErr        bool
	}{
		{
			"SuccessfulRequest",
			projectsService,
			args{
				UserId:    100,
				ProjectId: "default",
				RoleId:    3,
			},
			httpResponseMocks["success"],
			&AddUserRoleResponse{
				UserId:    100,
				ProjectId: "default",
				RoleId:    3,
			},
			&Response{Response: httpResponseMocks["success"]},
			false,
		},
		{
			"ReturnsError",
			projectsService,
			args{
				UserId:    100,
				ProjectId: "default",
				RoleId:    3,
			},
			httpResponseMocks["badrequest"],
			nil,
			&Response{Response: httpResponseMocks["badrequest"]},
			true,
		},
	}
	for _, tt := range tests {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return tt.mockedResponse, nil
		}
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := tt.p.DeleteUserProject(tt.args.ProjectId, tt.args.UserId, tt.args.RoleId)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectsService.DeleteUserProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// no additional checks for now as the response has no body.
			// ref: https://docs.getunleash.io/reference/api/legacy/unleash/admin/features-v2#remove-a-user-from-a-project
		})
	}
}
