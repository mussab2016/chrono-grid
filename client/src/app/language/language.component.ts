import {Component, inject} from '@angular/core';

import {MatMenuModule} from '@angular/material/menu';
import {MatButtonModule} from '@angular/material/button';
import {DEFAULT_SUPPORTED_LANGUAGES, Language, LanguageService} from "../services/languageService";

@Component({
    selector: 'language-menu',
    templateUrl: './language.component.html',
    styleUrls: ['./language.component.scss'], // <-- fixed typo
    standalone: true,
    imports: [
    MatMenuModule,
    MatButtonModule
]
})
export class LanguageComponent {
    public languages = DEFAULT_SUPPORTED_LANGUAGES;
    public languageSrv = inject(LanguageService);

    click(language: Language): void {
        this.languageSrv.setLanguage(language);
    }
}
