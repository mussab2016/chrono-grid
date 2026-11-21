import {Component, inject, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {CommonModule} from '@angular/common';
import {ErpNavTree, SideNavNode} from "../side-nav/side-nav.tree";
import {TranslateDirective} from "@common/translate-directive/translate.directive";
import {ShellService} from "@common/shell-service/shell.service";

@Component({
    selector: 'Main',
    standalone: true,
    imports: [CommonModule, TranslateDirective],
    templateUrl: './main.component.html',
    styleUrl: './main.component.scss'
})
export class MainComponent implements OnInit, OnDestroy {
    protected router = inject(Router);
    protected route = inject(ActivatedRoute);
    protected shellSrv = inject(ShellService);

    navNodes: SideNavNode[] = ErpNavTree;

    goToRoute(route: string) {
        this.router.navigate([route]).then();
    }

    getCircleColor(index: number): string {
        const colors = ['#e74c3c', '#3498db', '#2ecc71', '#f1c40f']; // example colors
        return colors[index % colors.length];
    }

    ngOnInit() {
        this.shellSrv.setTitle('unisched')
    }

    ngOnDestroy() {
        this.shellSrv.setTitle('')
    }

}

