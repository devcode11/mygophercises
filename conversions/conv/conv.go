package conv

import "fmt"

func DecToBase(num, base uint) (string, error) {
  
  if base > 36 {
    return "", fmt.Errorf("Base %d is not supported", base)
  }

  var result string

  for num > 0 {
    rem := num % base
    num = num / base
    if rem > 9 {
      rem += 65 - 10
    } else {
      rem += 48
    }
    result = string(rem) + result
  }
  return result, nil
}

func BaseToDec(num string, base uint) (uint, error) {
  if base > 36 {
    return 0, fmt.Errorf("Base %v not supported", base)
  }
  var result uint
  for _, digit := range num {
    var base10 uint
    if digit >= 65 {
      base10 = uint(digit - 65 + 10)
    } else {
      base10 = uint(digit - 48)
    }
    result = result * base + (base10)
  }
  return result, nil
}

// func ToBase(num string, fromBase, toBase uint) (string, error) {
//   if fromBase > 36 || toBase > 36 {
//     return "", fmt.Errorf("Maximum support base is 36")
//   }

  
// }