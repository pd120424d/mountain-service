import { HttpClientTestingModule, HttpTestingController } from "@angular/common/http/testing";
import { VersionInfo, VersionService } from "./version.service";
import { TestBed } from "@angular/core/testing";
import { environment } from "../../environments/environment";

describe('VersionService', () => {
  let service: VersionService;
  let httpMock: HttpTestingController;
  let expectedBaseUrl: string;
  let expectedVersionUrl: string;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [VersionService]
    });
    service = TestBed.inject(VersionService);
    httpMock = TestBed.inject(HttpTestingController);

    expectedBaseUrl = environment.useMockApi ? '/api/v1' : environment.apiUrl;
    expectedVersionUrl = expectedBaseUrl + '/version';
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should fetch version info', () => {
    const mockVersionInfo: VersionInfo = {
      version: '1.0.0',
      gitSHA: 'abc123',
      uptime: '1 hour'
    };

    service.getVersion().subscribe(versionInfo => {
      expect(versionInfo).toEqual(mockVersionInfo);
    });

    const req = httpMock.expectOne(expectedVersionUrl);
    expect(req.request.method).toBe('GET');
    req.flush(mockVersionInfo);
  });
});