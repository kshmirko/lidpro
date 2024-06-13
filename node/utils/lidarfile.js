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
      i+=1
    }
  
    ret.push(inpArray[i]);
    
    const retArr = new Uint8Array(ret.length)
  
    for (let i = 0; i < ret.length; i++){
      retArr[i] = ret[i]
    }
    
    return retArr.buffer
  }
  
  /**
   * 
   * @param {*} file_buffer 
   * @param {*} step 
   * @returns 
   */
  async function readLidarFile(file_buffer, step) {
    // для начала выполним фильтрацию
  
    const fileLen = file_buffer.byteLength - 18;
    const log2_filelen = Math.log2(fileLen)
  
    let inp_buffer = file_buffer
    if (log2_filelen != Math.floor(log2_filelen)) {
      inp_buffer = await filterBuffer(file_buffer)
    }
  
    const wordArray = new Uint16Array(inp_buffer)
  
    const profile = wordArray.slice(9,wordArray.byteLength)
    
    return {
      ProfileDT: new Date(wordArray[0], wordArray[1], wordArray[2], wordArray[3], wordArray[4], wordArray[5]),
      ProfLen: wordArray[6],
      Count: wordArray[7],
      RepRate: wordArray[8],
      Step: step,
      Data: Array.from(profile),
      ProfType: ''
    }
  }
  
  /**
   * Объединяет данные каналов DAT b DAK за одно и то же время
   * @param {Object} DAT 
   * @param {Object} DAK 
   * @returns List<Object> - combined channel
   */
  async function combineChannels(DAT, DAK) {
    let ret = []
  
    for (let dat of DAT) {
      for (let dak of DAK) {
        if ((dat.ProfileDT.toISOString() == dak.ProfileDT.toISOString())) {
          ret.push({
            ProfileDT: dat.ProfileDT,
            ProfLen: dat.ProfLen,
            Count: dat.Count,
            RepRate: dat.RepRate,
            Step: dat.Step,
            Dat: dat.Data,
            Dak: dak.Data
          })
        }
      }
    }
  
    ret.sort((a,b)=>a.ProfileDT-b.ProfileDT)
  
    return ret
  }

/**
 * Читает буфер зипованого файла с даннми измерений, распаковывет содержимое 
 * и выполняет усреднение данных
 * @param {Buffer} filename - uploaded zip file 
 * @param {Number} accumtime - accumulation time in seconds
 * @returns {Object} Структура 
 */
async function readLidarZipBuffer(filename, accumtime){

    // TODO: провестм рефракторинг для использования памяти для хранения загруженных файлов
    const zip = new StreamZip.async({
        file:filename,
    })

    let entriesCount = await zip.entriesCount;
    console.log(`Файлов в архиве: ${entriesCount}`);
    let entries = await zip.entries();

    for (const entry of Object.values(entries)) {

        const entryData = await zip.entryData(entry.name) 
        const data_i = await readLidarFile(entryData, accumtime)
    
        // data_i.ProfType = entry.name.split('.').pop()
        // if(data_i.ProfType.endsWith('dat')){
        //   DAT.push(data_i)
        // }else if(data_i.ProfType.endsWith('dak')){
        //   DAK.push(data_i)
        // }
        
      }
    
    let data = null
    return data
}

module.exports = {readLidarZipBuffer}