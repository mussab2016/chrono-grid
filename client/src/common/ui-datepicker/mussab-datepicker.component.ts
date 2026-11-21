// mussab-datepicker.component.ts
import {Component, Input, Optional, Self} from '@angular/core';
import {ControlValueAccessor, NgControl} from '@angular/forms';
import {MatFormFieldControl} from '@angular/material/form-field';
import {Subject} from 'rxjs';

import {MatInputModule} from '@angular/material/input';
import {MatDatepickerModule} from '@angular/material/datepicker';
import {MatNativeDateModule} from '@angular/material/core';


@Component({
    selector: 'mussab-datepicker',
    templateUrl: './mussab-datepicker.component.html',
    styleUrls: ['./mussab-datepicker.component.scss'],
    standalone: true,
    imports: [
    MatInputModule,
    MatDatepickerModule,
    MatNativeDateModule
],
    providers: [
        {provide: MatFormFieldControl, useExisting: MussabDatepickerComponent}
    ]
})
export class MussabDatepickerComponent
    implements ControlValueAccessor, MatFormFieldControl<Date> {

    static nextId = 0;
    readonly stateChanges = new Subject<void>();
    readonly id = `mussab-datepicker-${MussabDatepickerComponent.nextId++}`;

    @Input() placeholder: string = '';
    @Input() required = false;
    @Input() disabled = false;
    @Input() type: 'date' | 'datetime' | 'time' = 'date';

    value: Date | null = null;

    constructor(@Optional() @Self() public ngControl: NgControl) {
        if (this.ngControl != null) {
            this.ngControl.valueAccessor = this;
        }
    }

    // === ControlValueAccessor methods (required) ===
    writeValue(value: Date | null): void {
        this.value = value;
        this.stateChanges.next();
    }

    registerOnChange(fn: (value: any) => void): void {
        this.onChange = fn;
    }

    registerOnTouched(fn: () => void): void {
        this.onTouched = fn;
    }

    setDisabledState(isDisabled: boolean): void {
        this.disabled = isDisabled;
    }

    // Placeholders for callbacks
    onChange: (value: any) => void = () => {
    };
    onTouched: () => void = () => {
    };

    // === MatFormFieldControl required properties ===
    focused = false;
    errorState = false;
    controlType = 'mussab-datepicker';

    get empty() {
        return !this.value;
    }

    get shouldLabelFloat() {
        return !!this.value;
    }

    setDescribedByIds(_: string[]): void {
    }

    onContainerClick(_: MouseEvent): void {
    }

    // Custom handler when date changes
    onDateChange(date: Date | null) {
        this.value = date;
        this.onChange(date);
        this.stateChanges.next();
    }
}
