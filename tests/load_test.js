// import grpc from 'k6/net/grpc';
// import { check } from 'k6';
// import crypto from 'k6/crypto';

// const client = new grpc.Client();
// client.load(['proto'], 'ratelimiter.proto');

// const routes = ['login', 'signup', 'post'];
// // const routes = ['post'];

// // export const options = {
// //     vus: 1,
// //     iterations : 10,
// //     // duration: '30s',
// // };
// export const options = {
//     vus: 40,
//     duration: '30s',
// };

// // export const options = {
// //     vus: 10,
// //     iterations: 10,
// //     // duration: '30s',
// // };

// export default function () {
//     if (__ITER === 0) {
//         client.connect('localhost:6000', {
//             plaintext: true,
//         });
//     }

//     // Rotate through all registered route policies.
//     const resourceName = routes[__ITER % routes.length];

//     const res = client.invoke(
//         'rate_limiter.RateLimiter/Check',
//         {
//             clientIp: '127.0.0.1',
//             resourceName: resourceName,
//         },
//         {
//             metadata: {
//                 token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJvd25lck5hbWUiOiJkdW1teSIsImlzcyI6InJhdGVfbGltaXRlciIsImV4cCI6NDk0NTY0OTIzMH0.BwVVGFXEt8CkrQyg12RwvqA28bW8HeH9JsaKbdNYqwY",
//                 req_id: 
//             },
//         },
        
//     );

//     // An RPC can succeed even when the rate limiter denies a request.
//     check(res, {
//         [`${resourceName}: RPC completed`]: (r) => r && r.status === grpc.StatusOK,
//     });
// }


import grpc from 'k6/net/grpc';
import { check } from 'k6';

const client = new grpc.Client();
client.load(['proto'], 'ratelimiter.proto');

const routes = ['login', 'signup', 'post'];

function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
        const r = Math.random() * 16 | 0;
        const v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

export const options = {
    vus: 20,
    duration: '30s',
};

export default function () {
    if (__ITER === 0) {
        client.connect(__ENV.GRPC_ADDR, {
            plaintext: true,
        });
    }
    // if (__ITER === 0) {
    //     client.connect('localhost:6000', {
    //         plaintext: true,
    //     });
    // }

    const resourceName = routes[__ITER % routes.length];

    // Generate a unique ID for this RPC.
    const reqId = generateUUID();

    // Measure the complete synchronous RPC duration.
    const startTime = Date.now();

    const res = client.invoke(
        'rate_limiter.RateLimiter/Check',
        {
            clientIp: '127.0.0.1',
            resourceName: resourceName,
        },
        {
            metadata: {
                token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJvd25lck5hbWUiOiJkdW1teSIsImlzcyI6InJhdGVfbGltaXRlciIsImV4cCI6NDk0NTY0OTIzMH0.BwVVGFXEt8CkrQyg12RwvqA28bW8HeH9JsaKbdNYqwY",
                req_id: reqId,
            },
        },
    );

    const timeTaken = Date.now() - startTime;

    // Log only requests taking more than 10 seconds.
    if (timeTaken > 80) {
        console.log(`${reqId} : ${timeTaken}ms`);
    }

    check(res, {
        [`${resourceName}: RPC completed`]: (r) =>
            r && r.status === grpc.StatusOK,
    });
}
