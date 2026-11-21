import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
    name: 'unixDate',
    standalone: true,
    pure: true
})

export class UnixDatePipe implements PipeTransform {
    transform(value: number): string {
        const date = new Date(value * 1000);
        return `${date.getFullYear()}/${date.getMonth() > 8 ? date.getMonth() + 1 : '0' + (date.getMonth() + 1)}/${date.getDate() > 9 ? date.getDate() : '0' + date.getDate()}`;
    }
}

@Pipe({
    name: 'unixTime',
    standalone: true,
    pure: true
})

export class UnixTimePipe implements PipeTransform {
    transform(value: number): string {
        const date = new Date(value * 1000);
        return `${date.getHours() > 9 ? date.getHours() : '0' + date.getHours()} : ${date.getMinutes() > 9 ? date.getMinutes() : '0' + date.getMinutes()} : ${date.getSeconds() > 9 ? date.getSeconds() : '0' + date.getSeconds()}`;
    }
}
