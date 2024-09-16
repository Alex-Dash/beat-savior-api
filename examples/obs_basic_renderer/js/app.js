ref = {
    "type": "newsong",
    "song": {
        "songDataType": 1,
        "playerID": "76561198147664854",
        "songID": "26F8618BB3B4381F70602BEB05F7EA048866A705",
        "songDifficulty": "easy",
        "songName": "Harumachi Clover",
        "songArtist": "Hanasaka Yui (CV: M.A.O)",
        "songMapper": "Kival Evan",
        "gameMode": "Standard",
        "songDifficultyRank": 1,
        "songSpeed": 1,
        "songStartTime": 0,
        "songDuration": 34.283,
        "songJumpDistance": 25.34,
        "trackers": {
            "hitTracker": {
                "leftNoteHit": 20,
                "rightNoteHit": 21,
                "bombHit": 0,
                "maxCombo": 41,
                "nbOfWallHit": 0,
                "miss": 0,
                "missedNotes": 0,
                "badCuts": 0,
                "leftMiss": 0,
                "leftBadCuts": 0,
                "rightMiss": 0,
                "rightBadCuts": 0
            },
            "accuracyTracker": {
                "accRight": 113.190475,
                "accLeft": 111.75,
                "averageAcc": 112.487808,
                "leftSpeed": 27.7666264,
                "rightSpeed": 27.53605,
                "averageSpeed": 27.6485271,
                "leftHighestSpeed": 38.37508,
                "rightHighestSpeed": 36.97145,
                "leftPreswing": 1.74691033,
                "rightPreswing": 1.60082841,
                "averagePreswing": 1.67208791,
                "leftPostswing": 0.869471848,
                "rightPostswing": 0.839229763,
                "averagePostswing": 0.853982,
                "leftTimeDependence": 0.161059231,
                "rightTimeDependence": 0.145906717,
                "averageTimeDependence": 0.1532982,
                "leftAverageCut": [70, 11.75, 30],
                "rightAverageCut": [69.95238, 13.2380953, 30],
                "averageCut": [69.97561, 12.5121956, 30],
                "gridAcc": [-1, 112.6, 112.818184, -1, 110.375, -1, -1, 113.5, -1, 113, 114, -1],
                "gridCut": [0, 10, 11, 0, 8, 0, 0, 8, 0, 2, 2, 0]
            },
            "scoreTracker": {
                "rawScore": 29757,
                "score": 29757,
                "personalBest": 29800,
                "rawRatio": 0.9764397,
                "modifiedRatio": 0.9764397,
                "personalBestRawRatio": 0.9778507,
                "personalBestModifiedRatio": 0.9778507,
                "modifiersMultiplier": 1
            },
            "winTracker": {
                "won": true,
                "rank": "SS",
                "endTime": 33.996685,
                "nbOfPause": 1
            },
            "distanceTracker": {
                "rightSaber": 198.384567,
                "leftSaber": 202.967072,
                "rightHand": 47.8242531,
                "leftHand": 47.68143
            },
            "scoreGraphTracker": {
                "graph": {
                    "10": 0.987236142,
                    "11": 0.9844957,
                    "12": 0.982167244,
                    "14": 0.981958032,
                    "15": 0.9804383,
                    "16": 0.9807671,
                    "17": 0.981254935,
                    "18": 0.981627643,
                    "19": 0.9817752,
                    "2": 0.9826087,
                    "20": 0.9828908,
                    "21": 0.9823701,
                    "22": 0.981530845,
                    "23": 0.980698049,
                    "24": 0.979455233,
                    "25": 0.977538764,
                    "26": 0.9768735,
                    "27": 0.9761166,
                    "28": 0.9753399,
                    "29": 0.9750409,
                    "3": 0.9884058,
                    "4": 0.9911111,
                    "5": 0.99209106,
                    "6": 0.9925135,
                    "7": 0.99465394,
                    "8": 0.992666364,
                    "9": 0.990440249
                }
            }
        },
        "playDate": "2024-09-16"
    }
}

const sleep = ms => new Promise(r => setTimeout(r, ms));

const SONG_RANKS = [
    ["UNKNOWN_0", "#000000"],
    ["Easy", "#008055"],
    ["UNKNOWN_2", "#000000"],
    ["Normal", "#1268a1"],
    ["UNKNOWN_4", "#000000"],
    ["Hard", "#bd5500"], 
    ["UNKNOWN_6", "#000000"],
    ["Expert", "#b52a1c"], 
    ["UNKNOWN_8", "#000000"], 
    ["Expert+", "#7646af"],
]

function powerCurve(x) {
    return Math.pow(x, 5)
}

async function recalcGraph(data = ref.song.trackers.scoreGraphTracker.graph) {
    k_sorted = Object.keys(data).map(e=>parseInt(e)).sort((a,b)=>a-b)
    k_min = k_sorted[0]
    k_max = k_sorted[k_sorted.length-1]

    el = document.getElementById("graph")
    ctx = el.getContext("2d")
    width = el.width
    height = el.height
    padding = 15
    ctx.clearRect(0, 0, width, height);

    ctx.font = "bold 18px sans-serif";

    ctx.strokeStyle = "rgba(200, 200, 200, 0.3)";
    ctx.fillStyle = "rgba(200, 200, 200, 0.3)";
    ctx.lineWidth  = 2;
    for (let h_perc of [0, 0.25, 0.5, 0.625, 0.75, 0.875, 0.9375, 1]){
        let y = (height-padding*2)*(1-powerCurve(h_perc))+padding
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(width, y);
        ctx.stroke();
        if(h_perc>=.5){
            ctx.fillText(`${h_perc*100}%`, width-padding-60, y)
        }
    }

    for(let dx = 0; dx<=k_max; dx+=5){
        let x = (width-padding*2)*((dx/k_max))+padding
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x, height);
        ctx.stroke();
    }

    ctx.strokeStyle = "rgba(200, 255, 255, 0.9)";
    ctx.lineWidth  = 4;
    ctx.beginPath();
    ctx.moveTo(padding, padding);
    for(let dx of k_sorted){
        dy = data[`${dx}`]
        let x = (width-padding*2)*((dx/k_max))+padding
        let y = (height-padding*2)*(1-powerCurve(dy))+padding
        ctx.lineTo(x, y);
        ctx.moveTo(x, y);
    }
    ctx.stroke();

}

async function populateData(data = ref) {
    
    document.getElementById("song_win_state").innerText = data.song.trackers.winTracker.won?"LEVEL CLEARED":"LEVEL FAILED"
    
    // song info
    document.getElementById("song_name").innerText = data.song.songName
    document.getElementById("song_artist").innerText = data.song.songArtist
    document.getElementById("song_mapper").innerText = data.song.songMapper
    document.getElementById("song_diff").innerText = SONG_RANKS[data.song.songDifficultyRank][0]
    document.getElementById("song_diff").style.cssText = `--diff_clr: ${SONG_RANKS[data.song.songDifficultyRank][1]}`

    // overall scores
    document.getElementById("rating").innerText = data.song.trackers.winTracker.rank
    document.getElementById("score").innerText = `${`${data.song.trackers.scoreTracker.score}`.replace(/(\d)(?=(\d{3})+$)/g, '$1 ')} (${(data.song.trackers.scoreTracker.modifiedRatio*100).toFixed(2)}%)`
    document.getElementById("max_combo").innerText = data.song.trackers.hitTracker.maxCombo
    document.getElementById("misses").innerText = `${(data.song.trackers.hitTracker.missedNotes === 0 && data.song.trackers.hitTracker.badCuts === 0)?"None (FC)":data.song.trackers.hitTracker.missedNotes+data.song.trackers.hitTracker.badCuts}`
    document.getElementById("pauses").innerText = data.song.trackers.winTracker.nbOfPause
    document.getElementById("pb_score").innerText = `${`${data.song.trackers.scoreTracker.personalBest}`.replace(/(\d)(?=(\d{3})+$)/g, '$1 ')} (${(data.song.trackers.scoreTracker.personalBestModifiedRatio*100).toFixed(2)}%)`

    // left and right hands acc
    for (let side of [["l","left","Left"], ["r","right","Right"]]){
        document.getElementById(`${side[0]}_avg_acc`).innerText = data.song.trackers.accuracyTracker[`acc${side[2]}`].toFixed(2)
        document.getElementById(`${side[0]}_avg_speed`).innerText = (data.song.trackers.accuracyTracker[`${side[1]}Speed`]*3.6).toFixed(2)
        document.getElementById(`${side[0]}_max_speed`).innerText = (data.song.trackers.accuracyTracker[`${side[1]}HighestSpeed`]*3.6).toFixed(2)
        document.getElementById(`${side[0]}_hand_dist`).innerText = (data.song.trackers.distanceTracker[`${side[1]}Hand`]).toFixed(2)
        document.getElementById(`${side[0]}_saber_dist`).innerText = (data.song.trackers.distanceTracker[`${side[1]}Saber`]).toFixed(2)

        suff = ["pre", "acc", "post"]
        for(let idx=0; idx<3; idx++){
            document.getElementById(`${side[0]}_${suff[idx]}`).innerText = data.song.trackers.accuracyTracker[`${side[1]}AverageCut`][idx].toFixed(2)
        }
    }

    // block acc
    for(let idx in data.song.trackers.accuracyTracker.gridAcc){
        hits = data.song.trackers.accuracyTracker.gridCut[idx]
        hit_acc = data.song.trackers.accuracyTracker.gridAcc[idx]

        document.getElementById(`block_${idx}`).innerHTML = (hits>0 && hit_acc>0)?`<div>Hits: ${hits}</div><div>Acc: ${hit_acc.toFixed(2)}</div>`:""
        c_max = 200
        if(hits>0 && hit_acc>0){
            c_r = powerCurve(1-(115/hit_acc))*c_max 
            c_g = powerCurve(115/hit_acc)*c_max 
            document.getElementById(`block_${idx}`).style.cssText = `--block_bg: rgba(${c_r},${c_g},0,0.5);`
        }
    }

}

async function showAnim() {
    document.getElementById("wrapper").style.opacity = 1
    document.getElementById("wrapper").style.transform = "scale(1)"
    await sleep(10000)
    document.getElementById("wrapper").style.opacity = 0
    document.getElementById("wrapper").style.transform = "scale(0)"
    await sleep(1000)
}

var p_interval
function wsReconn() {
    clearInterval(p_interval)
    ws = new WebSocket('ws://127.0.0.1:1337/ws')

    ws.onmessage = e=>{
    	obj = JSON.parse(e.data)
        if(obj.type == "ping"){
            ws.send(JSON.stringify({type:"pong"}))
        }
    	if(obj.type === "newsong"){
            populateData(obj)
            recalcGraph(obj.song.trackers.scoreGraphTracker.graph)
            showAnim() 
        }
    }
    ws.onerror = async (e) => {
        console.log(e)
        await sleep(5000)
        wsReconn()
    }
    ws.onopen=(e)=>{
      console.log("Opened")
    	ws.send(JSON.stringify({type:"ping"}))
    	p_interval = setInterval(()=>{
    		ws.send(JSON.stringify({type:"ping"}))
    	}, 45000)
    }
}
wsReconn()