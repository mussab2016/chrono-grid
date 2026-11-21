export abstract class PaginationModel {

    public abstract get page(): number;
    public abstract set page(p: number);

    public abstract get pageSize(): number;
    public abstract set pageSize(s: number);

    public abstract get elementCount(): number;
    public abstract set elementCount(c: number);

    public onChange: () => void = () => null;

    public static newInstance(conf?: { page: number, pageSize: number, elementCount: number }): PaginationModel {
        const c = conf == null
            ? {page: 1, pageSize: 5, elementCount: 0}
            : conf;
        return new ExtendedPaginationModel(c.page, c.pageSize, c.elementCount);
    }

    public abstract copyModel(that: PaginationModel): void;
}

export class ExtendedPaginationModel extends PaginationModel {

    public propagatedPage;
    public propagatedPageSize;
    public propagatedElementCount;

    constructor(public mPage: number, public mPageSize: number, public mElementCount: number) {
        super();
        this.propagatedPage = mPage;
        this.propagatedPageSize = this.mPageSize;
        this.propagatedElementCount = mElementCount;
    }

    public onPageChange: () => void = () => null;
    public onPageSizeChange: () => void = () => null;
    public onElementCountChange: () => void = () => null;

    public get page(): number {
        return this.mPage;
    }

    public set page(p: number) {
        if (this.mPage === p) {
            return;
        }
        this.mPage = p;
        if (this.onPageChange != null) {
            this.onPageChange();
        }
    }

    public get pageSize(): number {
        return this.mPageSize;
    }

    public set pageSize(s: number) {
        if (this.mPageSize === s) {
            return;
        }
        this.mPageSize = s;
        if (this.onPageSizeChange != null) {
            this.onPageSizeChange();
        }
    }

    public get elementCount(): number {
        return this.mElementCount;
    }

    public set elementCount(c: number) {
        if (this.mElementCount === c) {
            return;
        }
        this.mElementCount = c;
        this.propagatedElementCount = c;
        if (this.onElementCountChange != null) {
            this.onElementCountChange();
        }
    }

    public fireChangeEvent(): void {
        if (this.propagatedPage === this.mPage
            && this.propagatedPageSize === this.mPageSize
            && this.propagatedElementCount === this.mElementCount) {
            return;
        }
        this.propagatedElementCount = this.mElementCount;
        this.propagatedPageSize = this.mPageSize;
        this.propagatedPage = this.mPage;
        if (this.onChange != null) {
            this.onChange();
        }
    }

    public copyModel(thatModel: PaginationModel): void {
        if (thatModel == null) {
            return;
        }
        const that: ExtendedPaginationModel = (thatModel as ExtendedPaginationModel);
        this.mPage = that.mPage;
        this.mPageSize = that.mPageSize;
        this.mElementCount = that.mElementCount;
        this.propagatedPage = that.propagatedPage;
        this.propagatedPageSize = that.propagatedPageSize;
        this.propagatedElementCount = that.propagatedElementCount;
    }
}


