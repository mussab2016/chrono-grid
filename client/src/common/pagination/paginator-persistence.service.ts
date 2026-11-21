import {Injectable} from '@angular/core';
import {PaginationModel} from "./pagination-model";

@Injectable({
    providedIn: 'root'
})
export class PaginatorPersistenceService {
    private paginatorMap = new Map<string, PaginationModel>();
    private searchMap = new Map<string, string>();

    get(title: string): PaginationModel {
        return this.paginatorMap.get(title) || PaginationModel.newInstance();
    }

    set(title: string, model: PaginationModel, search: string) {
        this.paginatorMap.set(title, model);
        this.searchMap.set(title, search);
    }

    getSearch(title: string) {
        return this.searchMap.get(title) || '';
    }
}

