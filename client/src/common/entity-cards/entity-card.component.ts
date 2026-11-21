import { Component, input, output } from '@angular/core';
import { EntityCardSchema } from './entity-card.model';
import {TranslatePipe} from '../../app/main/pipes/unix-transtale.pipe';
import {MatCard} from '@angular/material/card';

@Component({
    selector: 'app-entity-card',
    standalone: true,
    imports: [TranslatePipe, MatCard],
    templateUrl: './entity-card.component.html',
    styleUrl: './entity-card.component.scss'
})
export class EntityCardComponent<T extends { id: number }> {
    // Inputs using Angular Signals
    public item = input.required<T>();
    public schema = input.required<EntityCardSchema<T>>();
    public isSelected = input<boolean>(false);

    // Outputs
    public selectCard = output<number>();
    public openDetail = output<number>();

    public getDisplayValue(key: keyof T): string {
        const val = this.item()[key];
        return val !== null && val !== undefined ? String(val) : '';
    }

    public getInitial(): string {
        const title = this.getDisplayValue(this.schema().titleKey);
        return title ? title.charAt(0).toUpperCase() : '#';
    }
}
