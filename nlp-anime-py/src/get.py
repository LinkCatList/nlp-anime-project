import pandas as pd
import numpy as np
import nltk
import string
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
import tensorflow as tf
nltk.download('punkt')

from keras.models import load_model
from googletrans import Translator

import json

# выгрузим все необходимое и добавим необходимые функции 

my_model = load_model('14_text_classifier.hdf5') # изменить путь!!!!!!

with open('vocab.json', 'r') as f:
    mini_vocab = json.load(f)


def tokenize_text(text):
  tokenized_text = nltk.word_tokenize(text)
  tokens = [i.lower() for i in tokenized_text if ( i not in string.punctuation )]
  return tokens 

def EncodeTokens(ListTokens, VocabTokens): # токены -> числа
  res = []
  res = [VocabTokens.get(word, VocabTokens['<UNKW>']) for word in ListTokens]
  return [VocabTokens['<START>']] + res


########## основная ветка предсказания

labels = np.array(['Action', 'Adventure', 'Comedy', 'Drama', 'Fantasy', 'Sci-Fi',
       'Music', 'Kids', 'Slice', 'Hentai'])

# вставляем запрос пользователя и переводим по необходимости
