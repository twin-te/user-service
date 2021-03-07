export function deepContaining(obj: any): any {
  if (typeof obj === 'function' || typeof obj === 'undefined' || obj === null)
    return obj
  else if (Array.isArray(obj))
    return expect.arrayContaining(obj.map((n) => deepContaining(n)))
  else if (obj instanceof RegExp) return expect.stringMatching(obj)
  else if (typeof obj === 'object') {
    const tmp = { ...obj }
    Object.keys(tmp).forEach((k: any) => {
      tmp[k] = deepContaining(tmp[k])
    })
    return expect.objectContaining(tmp)
  } else return obj
}
