"use strict"

const StreamZip = require('node-stream-zip');

/**
 * Принимает входной буфер и фильтрует его
 * @param {Buffer} file_buffer 
 * @returns UInt8Array
 */
async function filterBuffer(file_buffer) {
    const inpArray = new Uint8Array(file_buffer);

    let ret = [];

    let i = 0;
    while (i < inpArray.byteLength - 1) {
        if ((inpArray[i] == 0x0d) && (inpArray[i + 1] == 0x0a)) {
            ret.push(inpArray[i + 1]);
            i += 1;
        } else {
            ret.push(inpArray[i]);
        }
        i += 1
    }

    ret.push(inpArray[i]);

    const retArr = new Uint8Array(ret.length)

    for (let i = 0; i < ret.length; i++) {
        retArr[i] = ret[i]
    }

    return retArr.buffer
}

/**
 * 
 * @param {Buffer} file_buffer  
 * @returns {Object}
 */
async function readLidarFile(file_buffer) {
    // для начала выполним фильтрацию

    const fileLen = file_buffer.byteLength - 18;
    const log2_filelen = Math.log2(fileLen)

    let inp_buffer = file_buffer
    if (log2_filelen != Math.floor(log2_filelen)) {
        inp_buffer = await filterBuffer(file_buffer)
    }

    const wordArray = new Uint16Array(inp_buffer)

    const profile = wordArray.slice(9, wordArray.byteLength)

    let vec = new Uint32Array(wordArray[6])


    for (let i = 0; i < wordArray[7]; i++) {
        for (let k = 0; k < wordArray[6]; k++) {
            vec[k] += profile[(i * wordArray[6]) + k]
        }
    }

    return {
        ProfileDT: new Date(wordArray[0], wordArray[1], wordArray[2], wordArray[3], wordArray[4], wordArray[5]),
        ProfileStopDT: new Date(wordArray[0], wordArray[1], wordArray[2], wordArray[3], wordArray[4], wordArray[5])+10000,
        ProfLen: wordArray[6],
        Count: wordArray[7],
        RepRate: wordArray[8],
        Data: Array.from(vec)
    }
}

/**
 * Объединяет данные каналов DAT b DAK за одно и то же время
 * @param {Object} DAT 
 * @param {Object} DAK 
 * @returns {List<Object>} - combined channel
 */
function combineChannels(DAT, DAK) {
    let ret = []

    for (let dat of DAT) {
        for (let dak of DAK) {
            if ((dat.ProfileDT.toISOString() == dak.ProfileDT.toISOString())) {
                ret.push({
                    ProfTime: dat.ProfileDT,
                    ProfLen: dat.ProfLen,
                    ProfCnt: dat.Count,
                    RepRate: dat.RepRate,
                    Dat: dat.Data,
                    Dak: dak.Data
                })
            }
        }
    }

    ret.sort((a, b) => a.ProfTime - b.ProfTime)

    return ret
}


/**
 * Выполняет усреднение профилей по заданному временному окну
 * @param {List<Object>} data - список сырых измерений
 * @param {number} accumtime - время накопления в секундах
 * @returns {List<Object>} - список изверений с требуемым накоплением
 */
function makeAccumulation(data, accumtime) {
    accumtime = accumtime*1000
    const N = data.length
    let i = 0
    let j = 0
    let ret = []
    console.log(data[0].ProfTime, data[2].ProfTime)
    while (i < N) {
        let res = {}
        j = i + 1
        let avg_dat = [...data[i].Dat];
        let avg_dak = [...data[i].Dak];
        let cnt = data[i].ProfCnt
        
        while((j<N)&&((data[j].ProfTime - data[i].ProfTime)<accumtime)){
            console.log(data[j].ProfTime - data[i].ProfTime)
            for(let k=0; k<avg_dat.length; k++){
                avg_dat[k] = avg_dat[k]+data[j].Dat[k]
                avg_dak[k] = avg_dak[k]+data[j].Dak[k]
            }
            cnt+=data[j].ProfCnt
            j+=1
        }
        res.ProfTime = data[i].ProfTime
        res.StopDate = data[j-1].ProfTime
        res.ProfCnt = cnt*100
        res.ProfLen = data[i].ProfLen
        res.RepRate = data[i].RepRate
        res.Dat = avg_dat
        res.Dak = avg_dak
        res.ProfDataDat = JSON.stringify(res.Dat),
        res.ProfDataDak = JSON.stringify(res.Dak),
        ret.push(res)
        i=j
    }

    return ret
}


/**
 * Читает буфер зипованого файла с даннми измерений, распаковывет содержимое 
 * и выполняет усреднение данных
 * @param {Buffer} filename - uploaded zip file 
 * @param {Number} accumtime - accumulation time in seconds
 * @returns {Object} Структура 
 */
async function readLidarZipBuffer(filename, accumtime) {

    // TODO: провестм рефракторинг для использования памяти для хранения загруженных файлов
    const zip = new StreamZip.async({
        file: filename,
    })

    let entriesCount = await zip.entriesCount;
    console.log(`Файлов в архиве: ${entriesCount}`);
    let entries = await zip.entries();
    let DAT = [];
    let DAK = [];
    for (const entry of Object.values(entries)) {

        const ftype = entry.name.split('.').pop().toLowerCase();
        if (ftype === "dat" || ftype == "dak") {
            const entryData = await zip.entryData(entry);
            const data_i = await readLidarFile(entryData);

            if (ftype === 'dat') {
                DAT.push(data_i)
            } else if (ftype === 'dak') {
                DAK.push(data_i)
            }
        }
    }


    let data = combineChannels(DAT, DAK)
    // Make accumulation
    console.log("Данные: ", data.length)
    let data1 = makeAccumulation(data, accumtime)

    return data1
}

module.exports = { readLidarZipBuffer }