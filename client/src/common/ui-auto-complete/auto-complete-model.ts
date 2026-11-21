import { Observable, of } from 'rxjs';
import { map } from 'rxjs/operators';
import { BaseEntityService } from '../service-template/service-template.component';

export type Id = number;

export interface AutoCompleteModel<T> {
    valueOf: (a: T | null) => Id | null;
    toString: (a: T | null) => string;
    getById: (a: Id | null | undefined) => Observable<T | null>;
    filter: (a: string) => Observable<T[]>;
    onObjectSelected?: (a: T | null) => void;
    onValueChanged?: (a: Id | null) => void;
    options: T[];
    cache: T | null; // يمثل العنصر المختار حالياً لـ UiAutoCompleteComponent
    addEntry: (entry: string) => Observable<number>;
    isPerfectMatch(a: T | null, b: string): boolean;
}

export class DefaultAutoCompleteModel<T extends { id?: number }> implements AutoCompleteModel<T> {
    #searchLength = 49;
    #itemsMap = new Map<Id, T>();
    #displayFn: (item: T) => string;

    public options: T[] = [];
    public cache: T | null = null;

    public filterLoadedData = (x: T[]): T[] =>
        (x ?? []).filter(e => e != null && e[this.searchField] !== undefined && e[this.searchField] !== null);

    constructor(
        public entitySrv: BaseEntityService<T>,
        public searchField: Extract<keyof T, string> = 'description' as Extract<keyof T, string>,
        displayFn?: (item: T) => string,
        public idField: Extract<keyof T, string> = 'id' as Extract<keyof T, string>
    ) {
        this.#displayFn = displayFn ?? ((x: T) => String(x?.[this.searchField] ?? ''));
    }

    onValueChanged = (_e: Id | null) => {};
    onObjectSelected = (_e: T | null) => {};

    valueOf = (x: T | null): Id | null =>
        x ? (x[this.idField] as unknown as Id) : null;

    toString = (x: T | null): string =>
        x ? this.#displayFn(x).substring(0, 50) : '';

    getById = (id: Id | null | undefined): Observable<T | null> => {
        if (id === null || id === undefined || id <= 0) {
            return of(null);
        }
        const cached = this.#itemsMap.get(id);
        if (cached) {
            return of(cached);
        }
        return this.entitySrv.get(id).pipe(
            map(item => {
                if (item) {
                    this.#itemsMap.set(id, item);
                }
                return item;
            })
        );
    };

    filter = (filter: string): Observable<T[]> =>
        this.entitySrv.getList(0, this.#searchLength, filter, this.searchField).pipe(
            map(x => {
                const list = this.filterLoadedData(x.data);
                list.forEach(item => {
                    const id = this.valueOf(item);
                    if (id !== null) this.#itemsMap.set(id, item);
                });
                this.options = list;
                return list;
            })
        );

    isPerfectMatch = (obj: T | null, filter: string): boolean => {
        if (!filter || !obj) return false;
        const query = filter.trim().toUpperCase();
        const rawField = String(obj[this.searchField] ?? '').trim().toUpperCase();
        const formatted = this.toString(obj).trim().toUpperCase();
        return rawField === query || formatted === query;
    };

    clear(): void {
        this.options = [];
        this.cache = null;
        this.#itemsMap.clear();
    }

    addEntry(entry: string): Observable<number> {
        const trimmed = entry?.trim();
        if (!trimmed) {
            return of(0);
        }
        const entity = {
            [this.idField]: 0,
            [this.searchField]: trimmed
        } as unknown as T;

        return this.entitySrv.save(entity).pipe(
            map(res => (typeof res === 'number' ? res : (res as any)?.[this.idField] ?? 0))
        );
    }
}
