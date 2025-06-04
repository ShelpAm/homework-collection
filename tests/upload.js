import http from 'k6/http';
import { check } from 'k6';

// 初始化阶段：加载文件内容
// const homeworkFile = http.file(open('dummy_homework.txt', 'b'), 'homework.txt');
const homeworkFile = http.file(open('SIR.zip', 'b'), 'homework.txt');

const payloads = [
    {
        username: '杨扬骁',
        schoolId: '202326202015',
        assignmentName: '五个一',
        homework: homeworkFile,
    },
    {
        username: '竺羽翔',
        schoolId: '202326202002',
        assignmentName: '五个一',
        homework: homeworkFile,
    },
    {
        username: '刘志远',
        schoolId: '202326202001',
        assignmentName: '五个一',
        homework: homeworkFile,
    },
];

export const options = {
    vus: payloads.length,
    duration: '5s',
};

const url = 'http://localhost:8080/api/process-homework/';


export default function() {
    // __ITER is the iteration number, __VU is the VU number
    const idx = (__VU - 1 + __ITER) % payloads.length;
    const pl = payloads[idx];
    const res = http.post(url, pl);
    check(res, { 'status is 200': (r) => r.status === 200 });
}
