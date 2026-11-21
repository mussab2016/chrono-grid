import {FocusMonitor} from '@angular/cdk/a11y';
import {coerceBooleanProperty} from '@angular/cdk/coercion';
import {
    Component, computed,
    DoCheck,
    ElementRef,
    HostBinding,
    input,
    Input,
    OnDestroy,
    OnInit,
    Optional,
    Self, signal
} from '@angular/core';
import {ControlValueAccessor, UntypedFormControl, NgControl, FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MatAutocompleteSelectedEvent, MatAutocompleteTrigger, MatAutocomplete} from '@angular/material/autocomplete';
import {MatFormFieldControl} from '@angular/material/form-field';
import {of, Subject, Subscription} from 'rxjs';
import {debounceTime, distinctUntilChanged, filter, switchMap, tap} from 'rxjs/operators';
import {AutoCompleteModel, Id} from './auto-complete-model';
import {MatInput} from '@angular/material/input';
import {MatOption} from '@angular/material/core';
import {TranslateDirective} from "../translate-directive/translate.directive";

@Component({
    selector: 'ss-auto-complete',
    templateUrl: './auto-complete.component.html',
    styleUrls: ['./auto-complete.component.scss'],
    providers: [
        {provide: MatFormFieldControl, useExisting: UiAutoCompleteComponent}
    ],
    standalone: true,
    imports: [
        MatInput,
        FormsModule,
        MatAutocompleteTrigger,
        ReactiveFormsModule,
        MatAutocomplete,
        MatOption,
        TranslateDirective
    ]
})
export class UiAutoCompleteComponent implements ControlValueAccessor, MatFormFieldControl<Id>, DoCheck, OnInit, OnDestroy {
    #_inactive = false;
    notFound = signal(false);
    allowAdd = signal(false);
    canAdd = input(false);

    constructor(
        @Optional() @Self() public ngControl: NgControl,
        private fm: FocusMonitor,
        private elRef: ElementRef<HTMLElement>) {
        if (this.ngControl != null) {
            this.ngControl.valueAccessor = this;
        }
        fm.monitor(elRef.nativeElement, true).subscribe(origin => {
            if (this.focused && !origin) {
                this.onTouched();
            }
            this.isFocused = !!origin;
            this.stateChanges.next();
        });
    }

    get value(): Id {
        return this.selected || 0;
    }

    @Input()
    set value(selected: Id) {
        this.setValue(selected, true);
    }

    @Input()
    set placeholder(t: string) {
        this.pPlaceholder = t;
        this.stateChanges.next();
    }

    get placeholder(): string {
        return this.pPlaceholder || '';
    }

    get focused() {
        return this.isFocused;
    }

    get empty(): boolean {
        return !this.autoControl.value;
    }

    get shouldLabelFloat() {
        return this.focused || !this.empty;
    }

    get required(): boolean {
        return this.isRequired;
    }

    @Input()
    set required(b: boolean | string) {
        this.isRequired = b === true || b === 'true';
        this.stateChanges.next();
    }

    get disabled(): boolean {
        return !!this.isDisabled;
    }

    set disabled(value: boolean) {
        this.isDisabled = coerceBooleanProperty(value) || this.#_inactive;
        if (this.isDisabled) {
            this.autoControl.disable();
        } else {
            this.autoControl.enable();
        }
        this.stateChanges.next();
    }

    @Input()
    set inactive(b: boolean) {
        this.#_inactive = b;
    }

    static nextId = 0;

    readonly stateChanges = new Subject<void>();
    autoControl: UntypedFormControl = new UntypedFormControl();
    filterSub: Subscription | undefined;
    @HostBinding() readonly id = `ss-auto-complete-${UiAutoCompleteComponent.nextId++}`;
    isRequired = false;
    errorState = false;
    lastErrorState = false;
    pPlaceholder: string | undefined;
    readonly controlType: 'ssAutoComplete' | undefined;
    readonly autofilled = false;
    isFocused = false;
    isDisabled: boolean | undefined;
    selected: Id | undefined;
    loading = signal(false);
    isInputInvalid = computed(() => !this.loading() && !!this.autoControl.value && this.selected == null);

    counter = 0; // identifies the query, used to order queries and discard delayed result if they expire

    //Covert to signal
    @Input()
    model: AutoCompleteModel<any> | undefined;

    lastFilterValue: string | undefined;

    @Input('aria-describedby') userAriaDescribedBy: string | undefined;

    setAutoControlValue(a: Object | null): void {
        this.autoControl.setValue(a);
        this.lastFilterValue = this.model?.toString(a);
        this.stateChanges.next();
    }

    callModelSelectEvent(obj: any) {
        if (this.model && this.model.onObjectSelected != null) {
            this.model.onObjectSelected(this.model.cache);
        }
    }

    callModelValueChange(id: Id) {
        if (this.model && this.model.onValueChanged) {
            this.model.onValueChanged(id);
        }
    }

    setValue(v: Id, fireOnChange: boolean) {
        if (!!this.model && v !== null && v !== undefined && v > 0) {
            // possibilites:
            // the value is in cache
            // the value is in options
            // the value does not exist
            if (this.model.valueOf(this.model.cache) === v) {
                this.selected = v;
                // this.autoControl.setValue(this.model.cache);
                this.setAutoControlValue(this.model.cache);
                this.model.options = [this.model.cache];
                if (fireOnChange) {
                    this.onChange(this.selected);
                    this.callModelValueChange(this.selected);
                }
            } else if (this.displayValueFromOptions(v)) {
                if (fireOnChange && this.selected) {
                    this.onChange(this.selected);
                    this.callModelValueChange(this.selected);

                }
            } else {
                // this.autoControl.setValue(null);
                //this.setAutoControlValue(null);
                this.loadValue(v, () => {
                    if (fireOnChange && this.selected) {
                        this.onChange(this.selected);
                        this.callModelValueChange(this.selected);
                    }
                });
            }
        } else {
            this.selected = undefined;
            if (fireOnChange) {
                this.onChange(null);
                this.callModelValueChange(-1);
            }
        }
        this.stateChanges.next();
    }

    displayValueFromOptions(v: Id): boolean {
        if (this.model && this.model.options) {
            const e = this.model.options.filter(o => this.model!.valueOf(o) === v);
            if (e.length === 1) {
                this.model.cache = e[0];
                this.selected = v;
                this.model.options = e;
                // this.autoControl.setValue(e[0]);
                this.setAutoControlValue(e[0]);
                return true;
            }
        }
        return false;
    }

    loadValue(value: Id, runnable: () => void) {
        if (this.model) {
            this.loading.set(true);
            this.model
                .getById(value)
                .subscribe({
                    next: e => {
                        this.model!.cache = e;
                        this.model!.options = [e];
                        this.selected = value;
                        // this.autoControl.setValue(e);
                        this.setAutoControlValue(e);
                        runnable();
                        this.loading.set(false);
                    }, error: e => {
                        this.loading.set(false);
                        throw e;
                    }
                });
        } else {
            throw new Error('Auto complete model is null');
        }
    }

    setDescribedByIds(ids: string[]) {
    }

    onContainerClick(event: MouseEvent) {
        if ((event.target as Element).tagName.toLowerCase() != 'input' && this.elRef.nativeElement.querySelector('input')) {
            this.elRef.nativeElement.querySelector('input')!.focus();
        }
    }

    // ControlValueAccessor

    setDisabledState(isDisabled: boolean): void {
        this.disabled = isDisabled;
    }

    registerOnChange(fn: any): void {
        this.onChange = fn;
    }

    registerOnTouched(fn: any): void {
        this.onTouched = fn;
    }

    writeValue(id: any): void {
        this.setValue(id, false);
    }

    onChange = (a: any) => {
    };

    onTouched = () => {
    };

    // Life cycles (birth, spreading corruption in the land, death, hell)
    ngDoCheck(): void {
        if (this.ngControl) {
            this.errorState = (this.ngControl.invalid && this.ngControl.touched) || false;
            if (this.lastErrorState !== this.errorState) {
                this.stateChanges.next();
                this.lastErrorState = this.errorState;
            }
        }
    }

    ngOnInit() {
        if (this.model) {
            this.filterSub = this.autoControl.valueChanges.pipe(
                filter(txt => typeof ('') == typeof (txt)),
                debounceTime(300),
                filter(txt => txt !== this.model!.toString(this.model!.cache)), // break the loop.
                distinctUntilChanged((a, b) => {
                    return a === b && b === this.lastFilterValue;
                }),
                tap(f => {
                    //this.lastFilterValue = f;
                    this.model!.cache = null;
                    this.model!.options = [];
                    this.setValue(0, true);
                }),
                switchMap(txt => {
                    this.counter++;
                    const id = this.counter;
                    this.loading.set(true);
                    return txt == null ? of([]) : this.model!.filter(txt).pipe(
                        tap(options => {
                            const curText = this.autoControl.value?.toString();
                            this.notFound.set(curText.length > 3 && options.length === 0);
                            if (id == this.counter) {
                                this.model!.options = options;
                                if (options.length === 1) {
                                    if (this.model!.isPerfectMatch(options[0], txt)) {
                                        this.setValue(this.model!.valueOf(options[0])!, true);
                                        // this.callModelSelectEvent(options[0]);
                                    }
                                }
                            }
                        })
                    );
                })
            ).subscribe(
                {
                    next: () => this.loading.set(false),
                    error: err => {
                        this.model!.options = [];
                        this.loading.set(false);
                        throw err;
                    }
                }
            );
        }
    }

    ngOnDestroy() {
        this.filterSub?.unsubscribe();
        this.stateChanges?.complete();
        this.fm.stopMonitoring(this.elRef.nativeElement);
    }


    displayFn = (x: any) => {
        return this.model ? this.model.toString(x) : '';
    };

    optionSelected(event: MatAutocompleteSelectedEvent) {
        if (this.model) {
            const obj = this.model.valueOf(event.option.value);
            this.setValue(obj!, true);
            this.callModelSelectEvent(obj);
        }
    }

    clearInput(): void {
        this.setValue(-1, true);
        this.setAutoControlValue(null);
    }

    hasNoMatch(): boolean {
        if (this.model) {
            return !this.loading() && this.model.options && this.model.options.length == 0;
        }
        return false
    }

    addCurrentText() {
        if (this.model) {
            this.model.addEntry(this.autoControl.value.toString()).subscribe();
        }
    }
}
