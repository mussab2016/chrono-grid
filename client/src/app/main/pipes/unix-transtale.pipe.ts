import {inject, Pipe, PipeTransform} from '@angular/core';
import {LanguageService} from "../../services/languageService";

@Pipe({
    name: 'translate',
    standalone: true,
    pure: false
})

export class TranslatePipe implements PipeTransform {
    private languageServices = inject(LanguageService);

    transform(value: string): string {
        return this.languageServices.translate(value)
    }
}
