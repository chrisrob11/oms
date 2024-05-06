local app = 'oms';
local port = 8080;

// OMS Application Deployment
local omsDeployment = {
  apiVersion: 'apps/v1',
  kind: 'Deployment',
  metadata: {
    name: app,
    labels: { app: app },
  },
  spec: {
    replicas: 1,
    selector: { matchLabels: { app: app } },
    template: {
      metadata: { labels: { app: app } },
      spec: {
        containers: [
          {
            name: app,
            image: 'chrisrob1111/oms:latest',
            ports: [{ containerPort: port }],
            envFrom: [
              {
                secretRef: {
                  name: 'postgres-credentials',
                },
              },
            ],
            env: [
              {
                name: 'DATABASE_URL',
                value: 'host=postgres dbname=postgres user=$(username) password=$(password) sslmode=disable',
              },
            ],
          },
        ],
      },
    },
  },
};

// OMS Application Service
local omsService = {
  apiVersion: 'v1',
  kind: 'Service',
  metadata: {
    name: app,
  },
  spec: {
    selector: { app: app },
    ports: [
      {
        port: port,
        targetPort: port,
      },
    ],
  },
};

// PostgreSQL Credentials Secret
local postgresCredentials = {
  apiVersion: 'v1',
  kind: 'Secret',
  metadata: { name: 'postgres-credentials' },
  type: 'Opaque',
  data: {
    username: std.base64('user'),
    password: std.base64('password'),
  },
};

// PostgreSQL Deployment
local postgresDeployment = {
  apiVersion: 'apps/v1',
  kind: 'Deployment',
  metadata: { name: 'postgres' },
  spec: {
    replicas: 1,
    selector: { matchLabels: { app: 'postgres' } },
    template: {
      metadata: { labels: { app: 'postgres' } },
      spec: {
        containers: [
          {
            name: 'postgres',
            image: 'postgres:latest',
            ports: [{ containerPort: 5432 }],
            env: [
              { name: 'POSTGRES_USER', valueFrom: { secretKeyRef: { name: 'postgres-credentials', key: 'username' } } },
              { name: 'POSTGRES_PASSWORD', valueFrom: { secretKeyRef: { name: 'postgres-credentials', key: 'password' } } },
            ],
          },
        ],
      },
    },
  },
};

// PostgreSQL Service
local postgresService = {
  apiVersion: 'v1',
  kind: 'Service',
  metadata: { name: 'postgres' },
  spec: {
    type: 'ClusterIP',
    ports: [{ port: 5432, targetPort: 5432 }],
    selector: { app: 'postgres' },
  },
};

[
  omsDeployment,
  omsService,
  postgresCredentials,
  postgresDeployment,
  postgresService,
]

