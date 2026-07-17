import grpc from 'k6/net/grpc';
import { check } from 'k6';

const client = new grpc.Client();
client.load(['proto'], 'ratelimiter.proto');

export const options = {
    vus: 10,
    duration: '10s',
};

export default function () {
    if (__ITER == 0) {
        client.connect('localhost:6000', {
            plaintext: true,
        });
    }

    const res = client.invoke('rate_limiter.RateLimiter/Check', {
        clientIp: '127.0.0.1',
        resourceName: 'login',
    }, 
    {
        metadata: {
            "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJvd25lck5hbWUiOiJnYXRld2F5IiwiaXNzIjoicmF0ZV9saW1pdGVyIiwiZXhwIjo0OTM3MTg2MTgzfQ.m8gHDWAE_VshELtU4oTSSkv2cYWirokgNHE7nJITcks"
        }
    });

    console.log(JSON.stringify(res.message));

}