import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { AuthStore } from '@app/store';
import { ButtonModule } from 'primeng/button';
import { CheckboxModule } from 'primeng/checkbox';
import { InputGroupModule } from 'primeng/inputgroup';
import { InputGroupAddonModule } from 'primeng/inputgroupaddon';
import { InputTextModule } from 'primeng/inputtext';
import { PasswordModule } from 'primeng/password';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    InputGroupModule,
    InputGroupAddonModule,
    InputTextModule,
    CheckboxModule,
    ButtonModule,
    FormsModule,
    PasswordModule,
    RouterModule,
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
})
export class LoginComponent {
  public authStore = inject(AuthStore);

  public request = {
    username: '',
    password: '',
  };

  public login(): void {
    const loginRequest = {
      request: this.request,
    };
    this.authStore.login(loginRequest);
  }
}
