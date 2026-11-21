import {Component, inject} from '@angular/core';
import {RouterOutlet} from '@angular/router';

import {LanguageService} from "./services/languageService";
import {ShellComponent} from "@common/shell-service/shell.component";

@Component({
    selector: 'app-root',
    standalone: true,
    imports: [ShellComponent, RouterOutlet],
    templateUrl: './app.component.html',
    styleUrl: './app.component.scss'
})
export class AppComponent {
    public languageService = inject(LanguageService);



}


