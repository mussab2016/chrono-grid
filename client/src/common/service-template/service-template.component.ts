import {inject} from '@angular/core';
import {HttpParams} from '@angular/common/http';
import {Sort} from '@angular/material/sort';
import {Observable} from 'rxjs';
import {map} from 'rxjs/operators';

import {HttpService} from '../http-service/http.service';
import {PaginationModel} from '../pagination/pagination-model';

export interface HttpListResponse<T> {
    total: number;
    data: T[];
}

export abstract class BaseEntityService<T extends { id?: number }> {
    protected httpSrv = inject(HttpService);

    private _sortColumn: string | undefined = undefined;
    protected parentId: number = -1;

    protected constructor(
        protected url: string,
        protected parentUrl?: string
    ) {
    }

    public getPage(
        pageModel: PaginationModel,
        filter?: string,
        extraParams?: HttpParams
    ): Observable<T[]> {
        // Convert 1-based UI page index to 0-based backend page index
        const zeroBasedPage = Math.max(0, pageModel.page - 1);
        const activeSort = this._sortColumn || undefined;
        return this.getList(zeroBasedPage, pageModel.pageSize, filter, activeSort, extraParams).pipe(
            map((res) => {
                pageModel.elementCount = res.total;
                return res.data ?? [];
            })
        );
    }

    public getList(
        page: number,
        pageSize: number,
        filter?: string,
        sortBy?: string,
        extraParams?: HttpParams
    ): Observable<HttpListResponse<T>> {
        // Only pass clean non-empty filter and sort values
        const cleanFilter = filter?.trim() ? filter.trim() : undefined;
        const cleanSort = (this._sortColumn || sortBy) || undefined;


        return this.httpSrv.getData<T>(
            this.getUrl(),
            page,
            pageSize,
            cleanFilter,
            cleanSort,
            extraParams
        );
    }

    public get(id: number): Observable<T | null> {
        return this.httpSrv.get(this.getUrl() + id);
    }

    public getAll(filter?: string, sortBy?: string): Observable<HttpListResponse<T>> {
        // Only pass clean non-empty filter and sort values
        const cleanFilter = filter?.trim() ? filter.trim() : undefined;
        const cleanSort = (this._sortColumn || sortBy) || undefined;

        let httpParams: HttpParams | undefined = undefined;

        if (cleanFilter) {
            httpParams = (httpParams ?? new HttpParams()).set('filter', cleanFilter);
        }

        if (cleanSort) {
            httpParams = (httpParams ?? new HttpParams()).set('sort', cleanSort);
        }

        return this.httpSrv.getList<T>(this.getUrl(), httpParams);
    }

    public remove(id: number): Observable<{ count: number }> {
        return this.httpSrv.remove(this.getUrl() + id);
    }

    public add(entity: T): Observable<number> {
        return this.httpSrv
            .post(this.getUrl(), entity)
            .pipe(map((res: { id: number }) => res.id));
    }

    public save(entity: T): Observable<number> {
        return entity.id ? this.update(entity) : this.add(entity);
    }

    public update(entity: T): Observable<number> {
        return this.httpSrv
            .put(this.getUrl() + entity.id, entity)
            .pipe(map((res: { count: number }) => res.count));
    }

    protected getUrl(): string {
        if (this.parentUrl) {
            if (this.parentId == null || this.parentId < 0) {
                throw new Error('No valid parent id was given!');
            }
            const normalizedParent = this.parentUrl.endsWith('/') ? this.parentUrl : `${this.parentUrl}/`;
            return `${normalizedParent}${this.parentId}${this.url}`;
        }
        return this.url.endsWith('/') ? this.url : `${this.url}/`;
    }

    // Sets sorting only when a direction is active; clears it when sorting is reset
    public set sortColumn(s: Sort) {
        if (!s || !s.direction) {
            this._sortColumn = undefined;
            return;
        }
        this._sortColumn = (s.direction === 'desc' ? '-' : '+') + s.active;
    }

    public get sortColumnValue(): string | undefined {
        return this._sortColumn;
    }
}
