import { Component } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { of } from 'rxjs';

import { BaseEntityListComponent } from './base-entity-list-component';
import { BaseEntityService } from '../service-template/service-template.component';
import { ShellService } from '../shell-service/shell.service';
import { FancyDialogService } from '../fancy-dialog/fancy-dialog.service';

// 1. Create a concrete mock service for testing
class MockUserService {
    sortColumn: any;
    getPage() {
        return of([]);
    }
    remove(id: number) {
        return of(void 0);
    }
}

