import {inject, Injectable} from '@angular/core';
import {HttpResponse} from '@angular/common/http';
import {Observable, timeout} from 'rxjs';
import {HttpClient, HttpHeaders, HttpParams} from '@angular/common/http';
import {map} from 'rxjs/operators';
import {HttpListResponse} from "../service-template/service-template.component";

@Injectable({
    providedIn: 'root'
})
export class HttpService {
    public baseUrl: string = 'http://127.0.0.1:8080/api';
    public httpClient = inject(HttpClient);

    private buildUrl(url: string): string {
        if (!this.baseUrl || url == null) {
            throw new Error(`Invalid URL configuration: baseUrl=${this.baseUrl}, url=${url}`);
        }
        const cleanBase = this.baseUrl.replace(/\/+$/, '');
        const cleanPath = url.startsWith('/') ? url : `/${url}`;
        const fullUrl = `${cleanBase}${cleanPath}`;

        if (fullUrl.endsWith('NaN') || fullUrl.includes('undefined')) {
            throw new Error(`Malformed URI detected: ${fullUrl}`);
        }
        return fullUrl;
    }

    public getData<T>(
        url: string,
        page: number,
        size: number,
        filterStr?: string,
        sortBy?: string,
        params?: HttpParams,
        requestTimeout = 0
    ): Observable<HttpListResponse<T>> {
        let httpParams = (params ?? new HttpParams()).set('page', `${page}-${size}`);

        if (filterStr) {
            httpParams = httpParams.set('filter', filterStr);
        }
        if (sortBy) {
            httpParams = httpParams.set('sort', sortBy);
        }

        const req$ = this.httpClient.get<HttpListResponse<T>>(this.buildUrl(url), {
            params: httpParams
        });

        return requestTimeout > 0 ? req$.pipe(timeout(requestTimeout)) : req$;
    }

    public get<T>(url: string, httpParams?: HttpParams): Observable<T> {
        return this.httpClient.get<T>(this.buildUrl(url), {
            params: httpParams
        });
    }

    public getList<T>(url: string, httpParams?: HttpParams): Observable<HttpListResponse<T>> {
        return this.httpClient.get<HttpListResponse<T>>(this.buildUrl(url), {
            params: httpParams
        });
    }

    public getRawList<T>(url: string, httpParams?: HttpParams): Observable<T[]> {
        return this.httpClient.get<T[]>(this.buildUrl(url), {
            params: httpParams
        });
    }

    public getFullResponse<T = unknown>(url: string, httpParams?: HttpParams): Observable<HttpResponse<T>> {
        return this.httpClient.get<T>(this.buildUrl(url), {
            observe: 'response',
            params: httpParams
        });
    }

    public getBlobResponse(url: string, httpParams?: HttpParams): Observable<HttpResponse<Blob>> {
        return this.httpClient.get(this.buildUrl(url), {
            observe: 'response',
            params: httpParams,
            responseType: 'blob'
        });
    }

    public post<T = any>(url: string, data?: unknown): Observable<T> {
        return this.httpClient.post<T>(this.buildUrl(url), data);
    }

    public upload(url: string, data?: unknown): Observable<any> {
        return this.httpClient.post(this.buildUrl(url), data, {
            reportProgress: true,
            observe: 'events'
        });
    }

    public put<T = any>(url: string, data?: unknown): Observable<T> {
        return this.httpClient.put<T>(this.buildUrl(url), data);
    }

    // noinspection JSUnusedGlobalSymbols
    public removeById<T = any>(url: string, id: number): Observable<T> {
        const slash = url.endsWith('/') ? '' : '/';
        return this.remove<T>(`${url}${slash}${id}`);
    }

    public removeParams<T = any>(url: string, params: HttpParams): Observable<T> {
        return this.httpClient.delete<T>(this.buildUrl(url), { params });
    }

    public remove<T = any>(url: string, data?: unknown): Observable<T> {
        return this.httpClient.delete<T>(this.buildUrl(url), { body: data });
    }

    // noinspection JSUnusedGlobalSymbols
    public postBlob(url: string, data?: unknown, requestTimeout = 600000): Observable<ArrayBuffer> {
        const req$ = this.httpClient.post(this.buildUrl(url), data, {
            responseType: 'arraybuffer'
        });
        return requestTimeout > 0 ? req$.pipe(timeout(requestTimeout)) : req$;
    }

    public putBlob(url: string, data?: unknown): Observable<ArrayBuffer> {
        return this.httpClient.put(this.buildUrl(url), data, {
            responseType: 'arraybuffer'
        });
    }

}

