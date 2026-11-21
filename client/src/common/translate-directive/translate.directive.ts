import {Directive, ElementRef, OnInit, Renderer2} from '@angular/core';
import {BaseDirective} from './base-directive';
import {LanguageService} from "../../app/services/languageService";
import {takeUntil} from "rxjs";


@Directive({
    selector: '[translate]',
    standalone: true
})
export class TranslateDirective extends BaseDirective implements OnInit {

    constructor(protected languageService: LanguageService, protected override el: ElementRef, protected override renderer: Renderer2) {
        super(el, renderer);
    }

    override ngOnInit(): void {
        super.ngOnInit();
        this.languageService.loadInitialLanguage();
    }

    protected setup(): void {
        this.languageService.languageChange.pipe(takeUntil(this.destroy)).subscribe(
            () => {
                this.replaceText();
                this.replaceAttributes();
            }
        );
    }

    protected replaceText(): void {
        if (this.key) {
            this.setText(this.getValue(this.key));
        }
    }

    protected replaceAttributes(): void {
        if (this.attributes.length > 0) {
            const keys: string[] = this.getAttributesKeys();
            const data: any = this.languageService.translateAll(keys);
            this.setAttributes(data);
        }
    }

    protected getValue(key: string): string {
        return this.languageService.translate(key);
    }

}
