from locust import HttpUser, task, between


class CircuitBreakerCase(HttpUser):
    wait_time = between(0, 4)
    host = "http://127.0.0.1:8080"

    @task
    def process(self):
        self.client.get('/process')
