import http from 'k6/http';
import { check, sleep } from 'k6';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

// Chiến lược Load Test: Ramping Up (Tăng dần tải)
// Giống như video mô tả: Tăng số user từ từ để tìm điểm chết (Breaking Point)
export const options = {
  stages: [
    { duration: '30s', target: 20 },  // Giai đoạn 1: Warm up (20 users)
    { duration: '1m',  target: 100 }, // Giai đoạn 2: Tải cao (100 users - High Load)
    { duration: '30s', target: 200 }, // Giai đoạn 3: Stress Test (200 users - Cực hạn)
    { duration: '30s', target: 0 },   // Giai đoạn 4: Cool down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // SLA: 99% request phải dưới 1s
    http_req_failed: ['rate<0.01'],                 // Tỷ lệ lỗi phải dưới 1%
  },
};

const BASE_URL = 'http://136.110.7.108:8081';

// Tạo report HTML đẹp mắt để chụp ảnh cho vào CV
export function handleSummary(data) {
  return {
    "report.html": htmlReport(data),
  };
}

export default function () {
  // Giả lập Payload giống thực tế
  const payload = JSON.stringify({
    userID: "stress-test-user-" + __VU, // ID Unique cho mỗi Virtual User
    pickup: {
        latitude: 10.762622,
        longitude: 106.660172
    },
    destination: {
        latitude: 10.823099,
        longitude: 106.629664
    }
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  // Gửi request tính tiền (Heavy Logic)
  const res = http.post(`${BASE_URL}/trip/preview`, payload, params);

  // Kiểm tra kết quả (Assertions)
  check(res, {
    'status is 200': (r) => r.status === 200,
    'latency < 500ms': (r) => r.timings.duration < 500,
  });

  // Nghỉ giữa các request (Think Time) - Giả lập hành vi người thật
  // Random từ 0.5s đến 1.5s
  sleep(Math.random() * 1 + 0.5);
}
