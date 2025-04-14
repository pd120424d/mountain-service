import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

export interface VersionInfo {
    version: string;
    gitSha: string;
    uptime: string;
}

@Injectable({ providedIn: 'root' })
export class VersionService {
    private baseApiUrl = environment.useMockApi
        ? '/api/v1' // Mock server URL
        : `${environment.apiUrl}`; // Real API
    constructor(private http: HttpClient) { }

    getVersion(): Observable<VersionInfo> {
        return this.http.get<VersionInfo>(`${this.baseApiUrl}/version`);
    }
}
