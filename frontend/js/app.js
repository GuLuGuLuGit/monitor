const API_BASE = 'http://localhost:8080/api';
let healthChart = null;

// 更新服务器时间
function updateServerTime() {
    const now = new Date();
    document.getElementById('server-time').textContent = now.toLocaleTimeString('zh-CN');
}

// 获取仪表板数据
async function fetchDashboardData() {
    try {
        const response = await fetch(`${API_BASE}/dashboard`);
        const data = await response.json();
        updateDashboard(data);
    } catch (error) {
        console.error('Failed to fetch dashboard data:', error);
    }
}

// 更新仪表板
function updateDashboard(data) {
    // 更新状态信息
    if (data.status) {
        document.getElementById('channels-healthy').textContent = data.status.channels_healthy || 0;
        document.getElementById('channels-total').textContent = data.status.channels_total || 0;

        const percentage = data.status.channels_total > 0
            ? Math.round((data.status.channels_healthy / data.status.channels_total) * 100)
            : 0;
        document.getElementById('channels-percentage').textContent = percentage + '%';
    }

    // 更新健康状态
    if (data.health) {
        const indicator = document.getElementById('health-indicator');
        const healthText = document.getElementById('health-text');
        const statusDot = indicator.querySelector('.status-dot');

        healthText.textContent = data.health.overall_status === 'healthy' ? '健康' :
                                 data.health.overall_status === 'warning' ? '警告' : '错误';

        indicator.className = 'status-indicator status-' + data.health.overall_status;
        statusDot.style.background = data.health.overall_status === 'healthy' ? '#10b981' :
                                     data.health.overall_status === 'warning' ? '#f59e0b' : '#ef4444';
    }

    // 更新代理列表
    if (data.agents) {
        const agentsList = document.getElementById('agents-list');
        if (data.agents.length === 0) {
            agentsList.innerHTML = '<p class="loading">暂无代理</p>';
        } else {
            agentsList.innerHTML = data.agents.map(agent => `
                <div class="agent-item">
                    <h3>${agent.agent_name || agent.agent_id}</h3>
                    <p>ID: ${agent.agent_id}</p>
                    <p>状态: <span class="status-${agent.status}">${agent.status}</span></p>
                </div>
            `).join('');
        }
    }

    // 更新会话列表
    if (data.sessions) {
        const sessionsList = document.getElementById('sessions-list');
        if (data.sessions.length === 0) {
            sessionsList.innerHTML = '<p class="loading">暂无活动会话</p>';
        } else {
            sessionsList.innerHTML = data.sessions.map(session => `
                <div class="session-item">
                    <h3>${session.recipient}</h3>
                    <p>频道: ${session.channel}</p>
                    <p>状态: ${session.status}</p>
                    <p>最后活动: ${new Date(session.last_activity).toLocaleString('zh-CN')}</p>
                </div>
            `).join('');
        }
    }

    // 更新最后更新时间
    if (data.updated) {
        document.getElementById('last-update').textContent = new Date(data.updated).toLocaleString('zh-CN');
    }
}

// 获取系统信息
async function fetchSystemInfo() {
    try {
        const response = await fetch(`${API_BASE}/system/info`);
        const data = await response.json();
        document.getElementById('openclaw-version').textContent = data.openclaw_version || 'unknown';
        document.getElementById('db-status').textContent = data.database_status === 'connected' ? '已连接' : '未连接';
    } catch (error) {
        console.error('Failed to fetch system info:', error);
    }
}

// 获取历史数据并更新图表
async function fetchHistoryAndUpdateChart() {
    try {
        const response = await fetch(`${API_BASE}/status/history?hours=24`);
        const data = await response.json();
        updateChart(data);
    } catch (error) {
        console.error('Failed to fetch history:', error);
    }
}

// 更新图表
function updateChart(historyData) {
    const ctx = document.getElementById('health-chart').getContext('2d');

    // 准备数据
    const labels = historyData.map(item => {
        const date = new Date(item.timestamp);
        return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
    }).reverse();

    const healthyData = historyData.map(item => item.channels_healthy).reverse();
    const totalData = historyData.map(item => item.channels_total).reverse();

    // 销毁旧图表
    if (healthChart) {
        healthChart.destroy();
    }

    // 创建新图表
    healthChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [
                {
                    label: '健康频道',
                    data: healthyData,
                    borderColor: '#10b981',
                    backgroundColor: 'rgba(16, 185, 129, 0.1)',
                    tension: 0.4,
                    fill: true
                },
                {
                    label: '总频道',
                    data: totalData,
                    borderColor: '#667eea',
                    backgroundColor: 'rgba(102, 126, 234, 0.1)',
                    tension: 0.4,
                    fill: true
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'top',
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        stepSize: 1
                    }
                }
            }
        }
    });
}

// 手动刷新数据
async function refreshData() {
    const btn = document.getElementById('refresh-btn');
    btn.disabled = true;
    btn.textContent = '刷新中...';

    try {
        await fetch(`${API_BASE}/status/refresh`, { method: 'POST' });
        // 等待 2 秒后重新获取数据
        setTimeout(async () => {
            await fetchDashboardData();
            await fetchHistoryAndUpdateChart();
            btn.disabled = false;
            btn.textContent = '刷新数据';
        }, 2000);
    } catch (error) {
        console.error('Failed to refresh data:', error);
        btn.disabled = false;
        btn.textContent = '刷新数据';
    }
}

// 初始化
async function init() {
    updateServerTime();
    setInterval(updateServerTime, 1000);

    await fetchSystemInfo();
    await fetchDashboardData();
    await fetchHistoryAndUpdateChart();

    // 每 60 秒自动刷新
    setInterval(async () => {
        await fetchDashboardData();
        await fetchHistoryAndUpdateChart();
    }, 60000);

    // 绑定刷新按钮
    document.getElementById('refresh-btn').addEventListener('click', refreshData);
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', init);
