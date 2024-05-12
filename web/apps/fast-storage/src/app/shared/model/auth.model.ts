export interface LoginRequest {
  request: {
    username: string;
    password: string;
  };
}

export interface GetNewTokenRequest {
  request: {
    refreshToken: string;
  };
}

export interface LogoutRequest {
  request: {
    refreshToken: string;
  };
}

export interface RegisterRequest {
  request: {
    username: string;
    password: string;
    confirmPassword: string;
    email: string;
    firstName: string;
    lastName: string;
  };
}

export interface LoginResponse {
  accessToken: string;
  error: string;
  errorDescription: string;
  expiresIn: number;
  idToken: string;
  notBeforePolicy: number;
  refreshExpiresIn: number;
  refreshToken: string;
  scope: string;
  sessionState: string;
  tokenType: 'Bearer';
}

export interface UserInfoResponse {
  acr: string;
  active: boolean;
  allowedOrigins: string[];
  aud: string[];
  azp: string;
  clientId: string;
  emailVerified: boolean;
  exp: number;
  iat: number;
  iss: string;
  jti: string;
  preferredUsername: string;
  realmAccess: {
    roles: string[];
  };
  resourceAccess: {
    account: {
      roles: string[];
    };
    masterRealm: {
      roles: string[];
    };
  };
  scope: string;
  sessionState: string;
  sid: string;
  sub: string;
  tokenType: string;
  typ: string;
  username: string;
}
